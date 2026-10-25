#!/usr/bin/env python3
"""Synchronize public wire definitions from the official Python gencode.

No Python SDK or runtime dependencies are imported. A protobuf Python runtime
is required to read descriptors. Existing Go-only compatibility definitions
remain available; upstream definitions are authoritative for shared names.
"""
import ast
import re
import subprocess
from pathlib import Path
from google.protobuf import descriptor_pb2 as pb

ROOT = Path(__file__).resolve().parent.parent
SOURCE = ROOT / 'docs/xai-sdk-python/xai_sdk-1.21.0/src/xai_sdk/proto/v6'
TYPES = {1:'double', 2:'float', 3:'int64', 4:'uint64', 5:'int32', 6:'fixed64', 7:'fixed32', 8:'bool', 9:'string', 12:'bytes', 13:'uint32', 15:'sfixed32', 16:'sfixed64', 17:'sint32', 18:'sint64'}

def enum_text(enum, indent=''):
    lines = [f'{indent}enum {enum.name} {{']
    if len({v.number for v in enum.value}) < len(enum.value):
        lines.append(indent + '  option allow_alias = true;')
    for value in enum.value:
        lines.append(f'{indent}  {value.name} = {value.number};')
    for r in enum.reserved_range:
        lines.append(f'{indent}  reserved {r.start}' + (f' to {r.end}' if r.end != r.start else '') + ';')
    lines.append(indent + '}')
    return '\n'.join(lines)

def message_text(message, indent=''):
    lines = [f'{indent}message {message.name} {{']
    for r in message.reserved_range:
        lines.append(f'{indent}  reserved {r.start}' + (f' to {r.end-1}' if r.end != r.start+1 else '') + ';')
    for enum in message.enum_type:
        lines.append(enum_text(enum, indent + '  '))
    entries = {n.name:n for n in message.nested_type if n.options.map_entry}
    for child in message.nested_type:
        if not child.options.map_entry:
            lines.append(message_text(child, indent + '  '))
    def kind(field):
        return TYPES.get(field.type, field.type_name)
    def field_text(field, inner):
        entry = entries.get(field.type_name.rsplit('.', 1)[-1])
        if entry:
            value_type = f'map<{kind(entry.field[0])}, {kind(entry.field[1])}>'
        else:
            prefix = 'optional ' if field.proto3_optional else 'repeated ' if field.label == 3 else ''
            value_type = prefix + kind(field)
        options = ' [deprecated = true]' if field.options.deprecated else ''
        return f'{inner}{value_type} {field.name} = {field.number}{options};'
    for field in message.field:
        if not field.HasField('oneof_index') or field.proto3_optional:
            lines.append(field_text(field, indent + '  '))
    for index, oneof in enumerate(message.oneof_decl):
        fields = [f for f in message.field if f.HasField('oneof_index') and f.oneof_index == index and not f.proto3_optional]
        if fields:
            lines.append(f'{indent}  oneof {oneof.name} {{')
            lines.extend(field_text(f, indent + '    ') for f in fields)
            lines.append(indent + '  }')
    lines.append(indent + '}')
    return '\n'.join(lines)

def block(text, kind, name):
    match = re.search(r'^' + kind + r'\s+' + re.escape(name) + r'\s*\{', text, re.M)
    if not match:
        return None
    depth = 0
    offset = match.start()
    for line in text[match.start():].splitlines(True):
        code = line.split('//', 1)[0]
        depth += code.count('{') - code.count('}')
        offset += len(line)
        if depth == 0:
            return match.start(), offset
    raise ValueError('Unclosed proto definition: ' + name)

def replace(text, kind, definition, rendered):
    span = block(text, kind, definition.name)
    if span:
        return text[:span[0]] + rendered + '\n' + text[span[1]:]
    return text.rstrip() + '\n\n' + rendered + '\n'

for path in sorted(SOURCE.glob('*_pb2.py')):
    tree = ast.parse(path.read_text())
    statement = next(n for n in tree.body if isinstance(n, ast.Assign) and any(isinstance(t, ast.Name) and t.id == 'DESCRIPTOR' for t in n.targets))
    descriptor = pb.FileDescriptorProto.FromString(ast.literal_eval(statement.value.args[0]))
    target = ROOT / 'proto' / descriptor.name
    text = target.read_text()
    # Regenerate only declarations whose wire descriptors need correction.
    selected = {
        'chat.proto': {'GetCompletionsRequest','GetChatCompletionResponse','GetChatCompletionChunk','OutputFile','FileContent','AttachmentSearch','CompactContextRequest','CompactContextResponse','IncludeOption','InlineCitation'},
        'files.proto': {m.name for m in descriptor.message_type} | {e.name for e in descriptor.enum_type},
        'types.proto': {'ChunkConfiguration','BytesConfiguration'},
        'usage.proto': {'ServiceTier','ServerSideTool'},
        'batch.proto': {'ListBatchRequestMetadataRequest','Batch'},
        'image.proto': {'ImageUrlContent'},
        'video.proto': {'VideoUrlContent','AudioUrlContent'},
    }.get(target.name, set())
    if not selected:
        continue
    for dependency in descriptor.dependency:
        if f'import "{dependency}";' not in text:
            text = text.replace('\npackage ' + descriptor.package + ';', '\npackage ' + descriptor.package + f';\n\nimport "{dependency}";')
    for enum in descriptor.enum_type:
        if enum.name == 'Ordering' and target.name == 'files.proto':
            span = block(text, 'enum', 'Ordering')
            if span: text = text[:span[0]] + text[span[1]:]
            enum.name = 'FilesOrdering'
            selected.add('FilesOrdering')
        if enum.name in selected:
            rendered = enum_text(enum)
            if enum.name == 'FilesOrdering':
                rendered = rendered.replace('{', '{\n  option allow_alias = true;', 1).replace('\n}', '\n  FILES_ORDERING_ASC = 0;\n  FILES_ORDERING_DESC = 1;\n}')
            if enum.name == 'FilesSortBy':
                rendered = rendered.replace('{', '{\n  option allow_alias = true;', 1).replace('\n}', '\n  FILES_SORT_BY_NAME = 1;\n}')
            text = replace(text, 'enum', enum, rendered)
    for message in descriptor.message_type:
        if message.name in selected:
            rendered = message_text(message)
            if target.name == 'files.proto': rendered = rendered.replace('.xai_api.Ordering', '.xai_api.FilesOrdering')
            text = replace(text, 'message', message, rendered)
    for service in descriptor.service:
        span = block(text, 'service', service.name)
        if span:
            definition = text[span[0]:span[1]]
            for method in service.method:
                if not re.search(r'\brpc\s+' + method.name + r'\(', definition):
                    input_type = ('stream ' if method.client_streaming else '') + method.input_type
                    output_type = ('stream ' if method.server_streaming else '') + method.output_type
                    definition = definition.rstrip().removesuffix('}') + f'  rpc {method.name}({input_type}) returns ({output_type}) {{}}\n}}\n'
            text = text[:span[0]] + definition + text[span[1]:]
    target.write_text(text)
    print(target.relative_to(ROOT))

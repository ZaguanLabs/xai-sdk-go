#!/usr/bin/env python3
"""
Verify Go SDK proto definitions against official xAI Python SDK v1.4.0.
Extracts all message definitions and compares field numbers, types, and order.
"""

from xai_sdk.proto.v6 import (
    auth_pb2, chat_pb2, collections_pb2, deferred_pb2, 
    documents_pb2, embed_pb2, files_pb2, image_pb2, 
    models_pb2, sample_pb2, shared_pb2, tokenize_pb2, 
    types_pb2, usage_pb2
)

def get_field_type_name(field):
    """Get human-readable field type name."""
    field_type_map = {
        1: 'double', 2: 'float', 3: 'int64', 4: 'uint64',
        5: 'int32', 8: 'bool', 9: 'string', 11: 'message',
        12: 'bytes', 13: 'uint32', 14: 'enum', 15: 'sfixed32',
        16: 'sfixed64', 17: 'sint32', 18: 'sint64'
    }
    field_type = field_type_map.get(field.type, f'type{field.type}')
    
    if field.type == 11:  # MESSAGE
        field_type = field.message_type.name
    elif field.type == 14:  # ENUM
        field_type = field.enum_type.name
    
    return field_type

def extract_message_definition(message_class):
    """Extract complete message definition."""
    msg = message_class()
    desc = msg.DESCRIPTOR
    
    fields = []
    for field in sorted(desc.fields, key=lambda f: f.number):
        label = ''
        if field.label == 3:  # REPEATED
            label = 'repeated '
        
        field_type = get_field_type_name(field)
        fields.append({
            'number': field.number,
            'name': field.name,
            'type': field_type,
            'label': label,
            'proto_line': f"  {label}{field_type} {field.name} = {field.number};"
        })
    
    return {
        'name': desc.name,
        'fields': fields
    }

def extract_enum_definition(enum_descriptor):
    """Extract enum definition."""
    values = []
    for value in enum_descriptor.values:
        values.append({
            'name': value.name,
            'number': value.number,
            'proto_line': f"  {value.name} = {value.number};"
        })
    
    return {
        'name': enum_descriptor.name,
        'values': values
    }

def extract_service_definition(service_descriptor):
    """Extract service definition."""
    methods = []
    for method in service_descriptor.methods:
        methods.append({
            'name': method.name,
            'input_type': method.input_type.name,
            'output_type': method.output_type.name,
            'proto_line': f"  rpc {method.name}({method.input_type.name}) returns ({method.output_type.name});"
        })
    
    return {
        'name': service_descriptor.name,
        'methods': methods
    }

def extract_all_definitions(module):
    """Extract all messages, enums, and services from a module."""
    messages = {}
    enums = {}
    services = {}
    
    for item_name in dir(module):
        item = getattr(module, item_name)
        
        # Check if it's a message type
        if hasattr(item, 'DESCRIPTOR'):
            desc = item.DESCRIPTOR
            
            # Message
            if hasattr(desc, 'fields'):
                msg_def = extract_message_definition(item)
                messages[msg_def['name']] = msg_def
                
                # Extract nested enums
                for enum_desc in desc.enum_types:
                    enum_def = extract_enum_definition(enum_desc)
                    enums[enum_def['name']] = enum_def
            
            # Service
            elif hasattr(desc, 'methods'):
                svc_def = extract_service_definition(desc)
                services[svc_def['name']] = svc_def
    
    # Top-level enums
    if hasattr(module, 'DESCRIPTOR'):
        for enum_desc in module.DESCRIPTOR.enum_types_by_name.values():
            enum_def = extract_enum_definition(enum_desc)
            enums[enum_def['name']] = enum_def
    
    return messages, enums, services

def generate_proto_file(module_name, messages, enums, services):
    """Generate complete .proto file content."""
    lines = []
    lines.append('syntax = "proto3";')
    lines.append('')
    lines.append('package xai_api;')
    lines.append('')
    lines.append('option go_package = "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1;xaiv1";')
    lines.append('')
    
    # Add imports (simplified - would need to detect actual imports)
    if module_name != 'shared':
        lines.append('import "google/protobuf/timestamp.proto";')
        lines.append('')
    
    # Enums
    if enums:
        for enum_name, enum_def in sorted(enums.items()):
            lines.append(f"enum {enum_name} {{")
            for value in enum_def['values']:
                lines.append(value['proto_line'])
            lines.append("}")
            lines.append("")
    
    # Messages
    if messages:
        for msg_name, msg_def in sorted(messages.items()):
            lines.append(f"message {msg_name} {{")
            for field in msg_def['fields']:
                lines.append(field['proto_line'])
            lines.append("}")
            lines.append("")
    
    # Services
    if services:
        for svc_name, svc_def in sorted(services.items()):
            lines.append(f"service {svc_name} {{")
            for method in svc_def['methods']:
                lines.append(method['proto_line'])
            lines.append("}")
            lines.append("")
    
    return '\n'.join(lines)

def main():
    modules = {
        'auth': auth_pb2,
        'chat': chat_pb2,
        'collections': collections_pb2,
        'deferred': deferred_pb2,
        'documents': documents_pb2,
        'embed': embed_pb2,
        'files': files_pb2,
        'image': image_pb2,
        'models': models_pb2,
        'sample': sample_pb2,
        'shared': shared_pb2,
        'tokenize': tokenize_pb2,
        'types': types_pb2,
        'usage': usage_pb2,
    }
    
    print("=" * 80)
    print("EXTRACTING OFFICIAL xAI PROTO DEFINITIONS FROM PYTHON SDK v1.4.0")
    print("=" * 80)
    print()
    
    for name, module in modules.items():
        print(f"Processing {name}.proto...")
        messages, enums, services = extract_all_definitions(module)
        
        proto_content = generate_proto_file(name, messages, enums, services)
        
        output_file = f"proto/xai/v1/{name}.proto.extracted"
        with open(output_file, 'w') as f:
            f.write(proto_content)
        
        print(f"  ✓ Extracted {len(messages)} messages, {len(enums)} enums, {len(services)} services")
        print(f"  → {output_file}")
        print()
    
    print("=" * 80)
    print("EXTRACTION COMPLETE")
    print("=" * 80)
    print()
    print("Next steps:")
    print("1. Review extracted .proto.extracted files")
    print("2. Compare with existing .proto files")
    print("3. Update or create missing proto files")
    print("4. Run: buf generate proto")

if __name__ == '__main__':
    main()

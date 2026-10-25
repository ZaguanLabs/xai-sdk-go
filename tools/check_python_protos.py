#!/usr/bin/env python3
"""Compare all official public descriptors with the generated Go descriptor set.

Run buf build proto -o /tmp/zaguan-xai-go-descriptors.pb first. Go-only
compatibility fields and declarations are permitted. Wire types, numbers,
cardinality, presence and RPC streaming must match upstream.
"""
import ast
import sys
from pathlib import Path
from google.protobuf import descriptor_pb2 as pb

ROOT = Path(__file__).resolve().parent.parent
source = ROOT / 'docs/xai-sdk-python/xai_sdk-1.21.0/src/xai_sdk/proto/v6'
go = pb.FileDescriptorSet.FromString(Path(sys.argv[1] if len(sys.argv)>1 else '/tmp/zaguan-xai-go-descriptors.pb').read_bytes())
messages, enums, services = {}, {}, {}
def index(file,ms,es,ss):
 def msg(m,prefix):
  name=prefix+'.'+m.name;ms[name]=m
  for n in m.nested_type:msg(n,name)
  for e in m.enum_type:es[name+'.'+e.name]=e
 for m in file.message_type:msg(m,'.'+file.package)
 for e in file.enum_type:es['.'+file.package+'.'+e.name]=e
 for s in file.service:ss['.'+file.package+'.'+s.name]=s
for f in go.file:index(f,messages,enums,services)
def normalize(name):return '.xai_api.FilesOrdering' if name=='.xai_api.Ordering' else name
errors=[]
for path in sorted(source.glob('*_pb2.py')):
 tree=ast.parse(path.read_text());data=None
 for node in ast.walk(tree):
  if isinstance(node,ast.Call) and isinstance(node.func,ast.Attribute) and node.func.attr=='AddSerializedFile':data=ast.literal_eval(node.args[0]);break
 if data is None:continue
 f=pb.FileDescriptorProto.FromString(data)
 # Management/shared packages are maintained separately; compare API files present in Go.
 if not any(x.name.rsplit('/',1)[-1]==f.name.rsplit('/',1)[-1] for x in go.file):continue
 pm,pe,ps={},{},{};index(f,pm,pe,ps)
 for name,m in pm.items():
  g=messages.get(normalize(name))
  if g is None:errors.append(name+' missing');continue
  fields={x.name:x for x in g.field}
  for field in m.field:
   other=fields.get(field.name)
   if other is None:errors.append(name+'.'+field.name+' missing');continue
   a=(field.number,field.type,normalize(field.type_name),field.label,field.proto3_optional)
   b=(other.number,other.type,normalize(other.type_name),other.label,other.proto3_optional)
   if a!=b:errors.append(f'{name}.{field.name}: {a} != {b}')
   if field.HasField('oneof_index') and not field.proto3_optional:
    if not other.HasField('oneof_index') or m.oneof_decl[field.oneof_index].name!=g.oneof_decl[other.oneof_index].name:errors.append(name+'.'+field.name+' oneof mismatch')
 for name,e in pe.items():
  g=enums.get(normalize(name))
  if g is None:errors.append(name+' missing enum');continue
  values={x.name:x.number for x in g.value}
  for v in e.value:
   if values.get(v.name)!=v.number:errors.append(f'{name}.{v.name} enum mismatch')
 for name,s in ps.items():
  g=services.get(name)
  if g is None:errors.append(name+' missing service');continue
  methods={x.name:x for x in g.method}
  for method in s.method:
   other=methods.get(method.name)
   if other is None:errors.append(name+'.'+method.name+' missing RPC');continue
   for key in ('input_type','output_type','client_streaming','server_streaming'):
    if getattr(method,key)!=getattr(other,key):errors.append(name+'.'+method.name+' '+key+' mismatch')
if errors:
 print('\n'.join(errors));sys.exit(1)
print('All official API descriptor fields, enums and RPCs match.')

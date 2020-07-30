# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: netrounds_agent_untrusted_metrics.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='netrounds_agent_untrusted_metrics.proto',
  package='netrounds.agent.metrics',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=b'\n\'netrounds_agent_untrusted_metrics.proto\x12\x17netrounds.agent.metrics\"\xc3\x01\n\x07Metrics\x12\x16\n\x0emeasurement_id\x18\x01 \x01(\x04\x12\x16\n\x0e\x63onfig_version\x18\x02 \x01(\x04\x12\x11\n\ttimestamp\x18\x03 \x01(\x03\x12\x38\n\x04tags\x18\x04 \x03(\x0b\x32*.netrounds.agent.metrics.Metrics.TagsEntry\x12\x0e\n\x06values\x18\x05 \x01(\x0c\x1a+\n\tTagsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"O\n\x07Message\x12\x11\n\tsend_time\x18\x01 \x01(\x03\x12\x31\n\x07metrics\x18\x02 \x01(\x0b\x32 .netrounds.agent.metrics.Metricsb\x06proto3'
)




_METRICS_TAGSENTRY = _descriptor.Descriptor(
  name='TagsEntry',
  full_name='netrounds.agent.metrics.Metrics.TagsEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='netrounds.agent.metrics.Metrics.TagsEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='netrounds.agent.metrics.Metrics.TagsEntry.value', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=221,
  serialized_end=264,
)

_METRICS = _descriptor.Descriptor(
  name='Metrics',
  full_name='netrounds.agent.metrics.Metrics',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='measurement_id', full_name='netrounds.agent.metrics.Metrics.measurement_id', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='config_version', full_name='netrounds.agent.metrics.Metrics.config_version', index=1,
      number=2, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='netrounds.agent.metrics.Metrics.timestamp', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='tags', full_name='netrounds.agent.metrics.Metrics.tags', index=3,
      number=4, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='values', full_name='netrounds.agent.metrics.Metrics.values', index=4,
      number=5, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_METRICS_TAGSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=69,
  serialized_end=264,
)


_MESSAGE = _descriptor.Descriptor(
  name='Message',
  full_name='netrounds.agent.metrics.Message',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='send_time', full_name='netrounds.agent.metrics.Message.send_time', index=0,
      number=1, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='metrics', full_name='netrounds.agent.metrics.Message.metrics', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=266,
  serialized_end=345,
)

_METRICS_TAGSENTRY.containing_type = _METRICS
_METRICS.fields_by_name['tags'].message_type = _METRICS_TAGSENTRY
_MESSAGE.fields_by_name['metrics'].message_type = _METRICS
DESCRIPTOR.message_types_by_name['Metrics'] = _METRICS
DESCRIPTOR.message_types_by_name['Message'] = _MESSAGE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Metrics = _reflection.GeneratedProtocolMessageType('Metrics', (_message.Message,), {

  'TagsEntry' : _reflection.GeneratedProtocolMessageType('TagsEntry', (_message.Message,), {
    'DESCRIPTOR' : _METRICS_TAGSENTRY,
    '__module__' : 'netrounds_agent_untrusted_metrics_pb2'
    # @@protoc_insertion_point(class_scope:netrounds.agent.metrics.Metrics.TagsEntry)
    })
  ,
  'DESCRIPTOR' : _METRICS,
  '__module__' : 'netrounds_agent_untrusted_metrics_pb2'
  # @@protoc_insertion_point(class_scope:netrounds.agent.metrics.Metrics)
  })
_sym_db.RegisterMessage(Metrics)
_sym_db.RegisterMessage(Metrics.TagsEntry)

Message = _reflection.GeneratedProtocolMessageType('Message', (_message.Message,), {
  'DESCRIPTOR' : _MESSAGE,
  '__module__' : 'netrounds_agent_untrusted_metrics_pb2'
  # @@protoc_insertion_point(class_scope:netrounds.agent.metrics.Message)
  })
_sym_db.RegisterMessage(Message)


_METRICS_TAGSENTRY._options = None
# @@protoc_insertion_point(module_scope)

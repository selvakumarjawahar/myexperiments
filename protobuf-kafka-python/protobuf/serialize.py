
from generated.callexecuter_metrics_pb2 import Metrics
import numbers

def protobuf_serialize(data):
    proto_msg = Metrics()
    proto_msg.stream_id = data['stream_id']
    proto_msg.timestamp.FromMilliseconds(data['timestamp'])
    for key, val in data.items():
        if isinstance(val, numbers.Integral):
            proto_msg.values[key].int_val = val
        else:
            proto_msg.values[key].float_val = val
    return proto_msg.SerializeToString()

from generated.callexecuter_metrics_pb2 import Metrics


def protobuf_deserialize(data):
    metrics = Metrics()
    metrics.ParseFromString(data)
    return metrics

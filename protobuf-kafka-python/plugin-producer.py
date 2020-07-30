from confluent_kafka import KafkaException
from confluent_kafka import Producer as ConfluentProducer
from generated.http_plugin_pb2 import Metrics as HTTPMetrics
from generated.netrounds_agent_untrusted_metrics_pb2 import Message as PluginMessage
import random
import time

SimulatedTags = {
    "stream.id": "42",
    "stream.name": "http stream",
    "task.id": "10",
    "task.name": "http monitor task",
    "monitor.id": "monitor",
    "monitor.name": "http monitor",
    "plugin.name": "http"
}

SimulatedHTTPMetrics = {
}


def serializeHTTPMetrics():
    httpmsg = HTTPMetrics()
    httpmsg.connect_time_avg = 2.0
    httpmsg.first_byte_time_avg = 10.2
    httpmsg.response_time_min = 23.45
    httpmsg.response_time_avg = 21.56
    httpmsg.response_time_max = 12.45
    httpmsg.size_avg = 6.0
    httpmsg.speed_avg = 7
    httpmsg.es_timeout = 8
    httpmsg.es_response = 9
    httpmsg.es = 10
    return httpmsg.SerializeToString()


def serializePluginMetrics():
    pluginmsg = PluginMessage()
    pluginmsg.metrics.measurement_id = 42
    pluginmsg.metrics.config_version = 1
    pluginmsg.metrics.timestamp = 3000
    for key, val in SimulatedTags.items():
        pluginmsg.metrics.tags[key] = val
    pluginmsg.metrics.values = serializeHTTPMetrics()
    pluginmsg.send_time = 1000
    return pluginmsg.SerializeToString()


def pushKafka(kproducer, topic):
    msg = serializePluginMetrics()
    try:
        kproducer.produce(topic,
                          key=str(random.randint(1, 100)),
                          value=msg)
    except Exception:
        print(Exception.message)


if __name__ == "__main__":

    kproducer = ConfluentProducer({
        'bootstrap.servers': 'localhost:9092',
        # 'debug': 'topic,broker',
        'queue.buffering.max.ms': 100})
    while True:
        print("Pushing metrics")
        pushKafka(kproducer, "plugin.metrics")
        time.sleep(5)

import random


from confluent_kafka import KafkaException
from confluent_kafka import Producer as ConfluentProducer
from protobuf.serialize import protobuf_serialize as serialize

def compile_avro(data):
    print("compiling avro")
    return

class Producer(object):
    def __init__(self, schema=None, bootstrap_servers=None, topic=None, serializer=None):
        bootstrap_servers = "localhost:9092"

        if isinstance(bootstrap_servers, (list, tuple)):
            bootstrap_servers = ','.join(bootstrap_servers)

        self._bootstrap_servers = bootstrap_servers
        self._schema_file = schema
        self._topic = topic
        self._serializer = 'protobuf' if serializer is None else serializer

    def __del__(self):
        if hasattr(self, '_confluent_producer'):
            print('Destructor called, flushing Kafka producer')
            self.confluent_producer.flush()

    @property
    def confluent_producer(self):
        try:
            return self._confluent_producer
        except AttributeError:
            if self._bootstrap_servers:
                self._confluent_producer = ConfluentProducer({
                    'bootstrap.servers': self._bootstrap_servers,
                    # 'debug': 'topic,broker',
                    'queue.buffering.max.ms': 100})
                print('Created Kafka producer, bootstrap_servers=%s' % self._bootstrap_servers)
                return self._confluent_producer
            print('Kafka producer: None')
            return None

    def push(self, message, topic=None):
        """
        Push message to Kafka.

        :param message: Message to push
        """
        topic = topic or self._topic

        try:
            compiled = serialize(message) if self._serializer is 'protobuf' else self.compile_avro(message)
            self.confluent_producer.produce(topic,
                                            key=str(random.randint(1, 100)),
                                            value=compiled)
        except (ValueError, AttributeError):
            print('Failed to compile avro message')
        except KafkaException:
            print('Failed to push kafka message')
        except Exception:
            print('Unexpected error while pushing kafka message')

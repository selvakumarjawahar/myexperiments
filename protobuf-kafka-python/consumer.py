from confluent_kafka import Consumer
from confluent_kafka import KafkaError
from confluent_kafka import KafkaException
from protobuf.deserialize import protobuf_deserialize as deserialize
import sys

running = True
# Some Reason this Function not working
def consume_from_kafka(consumer, topics):

    try:
        consumer.subscribe(topics)
        while running:
            msg = consumer.poll(timeout=5.0)
            if msg is None:
                print('No message')
                continue

            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # End of partition event
                    sys.stderr.write('%% %s [%d] reached end at offset %d\n' %
                                     (msg.topic(), msg.partition(), msg.offset()))
                elif msg.error():
                    raise KafkaException(msg.error())
            else:
                print (deserialize(msg))

    finally:
        # Close down consumer to commit final offsets.
        consumer.close()

def consume_from_binary_file(filename):
    with open(filename, "rb") as f_bin:
        msg = f_bin.read()
        print (deserialize(msg))

if __name__ == "__main__":
    # Used the Following command to get the msg
    # kafkacat -C -b localhost:9092 -t test-proto -c 1 -f %s > msg.bin
    consume_from_binary_file('msg.bin')




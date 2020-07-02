import os
from kafka import Producer

test_data = {
    "timestamp": 2345670,
    "stream_id": 1,
    "rate": 5,
    "dmin": 1,
    "davg": 3,
    "dmax": 5,
    "dv": 23,
    "recv": 300,
    "loss_far": 1.2,
    "lost_far": 3.4,
    "miso_far": 23,
    "dmin_far": 23,
    "davg_far": 35,
    "dmax_far": 12,
    "dv_near": 23,
    "es": 1.23,
    "es_loss": 2.33,
    "es_delay": 45,
    "es_dv": 5.657,
    "es_dscp": 23,
    "es_ses": 24,
    "uas": 21
    }

if __name__ == "__main__":
    kafka_prod = Producer(topic='test-proto')
    print("Pushing Test Data")
    kafka_prod.push(test_data)
    print("Test Data pushed")


# Description
This package container sample code for subcribing to nats url/topic and publishing to nats url/topic.

## Layout
1. `demo-app.py`: Sample boiler plate code that a python app running on xi-iot would need in order to leverage input/output data sources
2. `xi-iot.proto`: Protobuf definition for serializing/deseriailizing incoming/outgoing messages
3. `publisher-example.py`: A helper python script to publish messages using the required protobufs(#2)

# Setup
1. The sample code features from python 3.5+. Hence, please update your python interpreter to an approriate version.
2. Install python nats client `pip3 install install asyncio-nats-client # install the python nats client`
3. You would need to have a protoc compiler installed in order to use the above mentioned protobufs.  https://developers.google.com/protocol-buffers/docs/downloads(Tested version: libprotoc 3.7.1)
4. With current working dir as this dir, run `protoc  --python_out=. xi-iot.proto # this generates xi_iot_pb2.py` 

import asyncio
import os
import xi_iot_pb2
from nats.aio.client import Client as NATS
from nats.aio.errors import ErrConnectionClosed, ErrTimeout, ErrNoServers

nc = None
nats_broker_url,src_nats_topic,dst_nats_topic = "","",""
msgs_received, msgs_sent = 0, 0

def get_nats_meta():
    global nats_broker_url,src_nats_topic,dst_nats_topic
    nats_broker_url = os.environ.get('NATS_ENDPOINT')
    if nats_broker_url is None:
        print('nats broker not provided in environment var NATS_ENDPOINT')
        exit(1)

    src_nats_topic = os.environ.get('NATS_SRC_TOPIC')
    if src_nats_topic is None:
        print('src nats topic not provided in environment var NATS_SRC_TOPIC')
        exit(1)
    dst_nats_topic = os.environ.get('NATS_DST_TOPIC')

    if dst_nats_topic is None:
        print('dst nats broker not provided in environment var NATS_DST_TOPIC')
        exit(1)
    return nats_broker_url, src_nats_topic, dst_nats_topic

async def message_handler(msg):
    # print('recieved msg')
    global msgs_received, msgs_sent
    subject = msg.subject
    reply = msg.reply
    # data = msg.data.decode()
    # print("Received a message on '{subject} {reply}': {data}".format(
    #    subject=subject, reply=reply))
    msgs_received += 1
    if msgs_received % 100 == 0:
      print('msgs_received: ', msgs_received)
    _msg = xi_iot_pb2.DataStreamMessage()
    _msg.ParseFromString(msg.data)
    # print("processed {data}".format(data=_msg.SerializeToString()))
    # ***************** your app's business logic here ********************
    # RFC: We could leverage `reply` topic as the destination topic which would not require DST_NATS_TOPIC to be provided
    # await nc.publish(reply, data)
    await nc.publish(dst_nats_topic, _msg.SerializeToString())
    msgs_sent += 1
    if msgs_sent % 100 == 0:
      print('msgs sent: ', msgs_sent)

async def run(loop):
    nats_broker_url, src_nats_topic, dst_nats_topic = get_nats_meta()
    print ("broker: {b}, src topic: {s}, dst_topic: {d}".format(b=nats_broker_url, s=src_nats_topic, d=dst_nats_topic))

    global nc
    nc = NATS()

    # This will return immediately if the server is not listening on the given URL
    await nc.connect(loop=loop, servers=[str(nats_broker_url)])
    print('connected')
    await nc.subscribe(str(src_nats_topic), cb=message_handler)

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(run(loop))
        loop.run_forever()
    finally:
        nc.drain()
        loop.close()

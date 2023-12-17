import anyio
import asyncio

import redis.asyncio as redis

STOPWORD = "STOP"


async def reader(channel: redis.client.PubSub):
    while True:
        # message = await channel.get_message(ignore_subscribe_messages=True)
        message = await channel.get_message(ignore_subscribe_messages=False)
        if message is not None:
            print(f"(Reader) Message Received: {message}")
            try:
                if message["data"].decode() == STOPWORD:
                    print("(Reader) STOP")
                    break
            except Exception as e:
                print(f"(Reader) Exception: {e}")

async def main():

    r = await redis.from_url("redis://localhost")
    async with r.pubsub() as pubsub:
        await pubsub.psubscribe("channel:*")

        future = asyncio.create_task(reader(pubsub))

        await r.publish("channel:1", "Hello")
        await r.publish("channel:2", "World")
        await r.publish("channel:3", "again")
        await r.publish("channel:1", STOPWORD)

        await future

anyio.run(main)
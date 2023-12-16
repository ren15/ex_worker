import os
from typing import Optional
import logging

import redis
from fastapi import FastAPI
import uvicorn

DEBUG=os.environ.get("DEBUG", False)
CHANNEL = os.environ.get("CHANNEL")
REDIS_HOST = os.environ.get("REDIS_HOST")


def publish(message):
    while True:  # note: limit this to x attempts, not a good idea to try indefinitely
        global r
        try:
            # if(random.randint(0,9) < 3):
            #     raise redis.ConnectionError("Test Connection Error")
            rcvd = r.publish(CHANNEL, message)
            if rcvd >0:
                break
        except redis.ConnectionError as e:
            logging.error(e)
            logging.error("Will attempt to retry")
        except Exception as e:
            logging.error(e)
            logging.error("Other exception")


app = FastAPI()
r = redis.Redis(host=REDIS_HOST)


@app.get("/")
async def root():
    return "OK"


@app.post("/messaging/send")
async def send_message(
    message: Optional[str] = ''
    ):
    if message != '':
        publish(message)
    else:
        publish("test message")
    return {"status": "succes"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
export CHANNEL=chan1
export REDIS_HOST=localhost
export PYTHONUNBUFFERED=1

python src/python/publisher.py &
# get pid
server_pid=$!

sleep 3

curl -X POST \
    -H 'Content-Type: text/plain' \
    --data "hello" \
    http://127.0.0.1:8000/messaging/send

kill $server_pid

echo "Server killed, exiting..."

export CHANNEL=chan1
export REDIS_HOST=localhost

python src/python/publisher.py &
# get pid
server_pid=$!

sleep 3

curl -X POST \
    http://127.0.0.1:8000/messaging/send \
    -H 'Content-Type: application/text' \
    "hello"

kill $server_pid

echo "Server killed, exiting..."

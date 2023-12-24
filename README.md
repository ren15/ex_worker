# ex_worker

This repo is demoing in elixir as a worker.

It subs from redis, and do it's work when some jobs are published from python.


## Progress
DONE:

1. start redis
2. python pub sub

TODO:

4. elixir sub 
5. run everything in github action.
6. elixir consumer from redpanda https://hexdocs.pm/broadway/apache-kafka.html#getting-started https://hexdocs.pm/broadway/architecture.html#the-supervision-tree

7. golang callback from sqlite
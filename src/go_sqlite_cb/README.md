
# go_sqlite_cb

## DONE
- golang callback from sqlite `src/go_sqlite_cb`
- callback call http
- callback send msg to global channel 

## TODO

- check callback latency
- two processes, use sqlite trigger callback to send msg to another process (maybe impossible because trigger is attached to a db connection, and the connection is not shared between processes)


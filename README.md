lagproxy
===

Lagproxy is a simple netcat-like TCP proxy that will simulate nasty network lag.

If your webserver is on localhost port 80, you could run lagproxy this way:

./lagproxy --from localhost:8080 --to localhost:80

To control how much lag happens, use --min, --max, and --block.  --block is the block size in bytes.  --min and --max control the lower and upper bounds (in milliseconds) to randomly sleep for each block read.

License
---

BSD License

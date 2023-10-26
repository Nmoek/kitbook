#!/usr/bin/env sh

echo "start signup test..."


ab -n 2 -c 1 -t 1 -T "application/json" -p http://localhost:80/users/signup

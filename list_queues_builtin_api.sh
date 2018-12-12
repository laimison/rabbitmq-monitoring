#!/bin/bash

echo "From guest account"
curl -i -u guest:guest http://localhost:15672/api/queues | tail -n 1 | jq

echo "From mon account"
curl -i -u mon:password123 http://localhost:15672/api/queues | tail -n 1 | jq

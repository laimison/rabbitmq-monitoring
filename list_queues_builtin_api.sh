#!/bin/bash

echo "From guest account"
curl -i -u guest:guest http://localhost:15672/api/queues | tail -n 1 | jq

echo "From monitoring account"
curl -i -u monitoring:password http://localhost:15672/api/queues | tail -n 1 | jq

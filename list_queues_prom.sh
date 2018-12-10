#!/bin/bash

curl -i -u guest:guest http://localhost:15672/api/queues | tail -n 1 | jq

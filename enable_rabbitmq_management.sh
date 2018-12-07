#!/bin/bash

rabbitmq-plugins enable rabbitmq_management
lsof -n -i -P | grep LISTEN | grep 15672
echo
curl http://localhost:15672

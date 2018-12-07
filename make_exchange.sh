#!/bin/bash

rabbitmqadmin declare exchange --vhost=Some_Virtual_Host name=some_exchange type=direct || exit 1
echo success

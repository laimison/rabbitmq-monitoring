#!/bin/bash

if ! brew services list | grep -q ^rabbitmq
then
  echo "Installing MQ"
  brew install rabbitmq
fi

if ! brew services list | grep ^rabbitmq | grep ' started '
then
  echo "Starting MQ"
  brew services start rabbitmq
fi

p='export PATH=$PATH:/usr/local/sbin'
if ! grep -q "$p" ~/.bash_profile
then
  echo "Adding rabbitmq PATH to profile"
  printf "\n$p\n" >> ~/.bash_profile
fi

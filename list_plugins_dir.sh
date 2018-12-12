#!/bin/bash

rabbitmqctl eval 'application:get_env(rabbit, plugins_dir).'

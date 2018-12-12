#!/bin/bash

# set plugins dir
# mkdir -p /usr/lib/rabbitmq/plugins
# cd /usr/lib/rabbitmq/plugins
plugins_dir=`rabbitmqctl eval 'application:get_env(rabbit, plugins_dir).' | awk -F '"' '{print $2}'`
cd $plugins_dir

# Downloads prometheus_rabbitmq_exporter and its dependencies with curl

readonly base_url='https://github.com/deadtrickster/prometheus_rabbitmq_exporter/releases/download/v3.7.2.3'

get() {
  curl -LO "$base_url/$1"
}

echo 'Press enter'
read

get accept-0.3.3.ez
get prometheus-3.5.1.ez
get prometheus_cowboy-0.1.4.ez
get prometheus_httpd-2.1.8.ez
get prometheus_process_collector-1.3.1.ez
get prometheus_rabbitmq_exporter-3.7.2.3.ez

rabbitmq-plugins enable prometheus_rabbitmq_exporter

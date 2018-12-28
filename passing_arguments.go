//usr/bin/env go run "$0" "$@"; exit

// An example to call this script (remove loop function)
//
// loop ./passing_arguments.go --queues-ignore ignore_this_queue --queues-ignore ignore_this_queue_as_well --warning-threshold 1 --critical-threshold 2 --warning-threshold 1 --critical-threshold 2 --default-warning-threshold 4 --default-critical-threshold 5 --queue some_incoming_queue --queue some_outgoing_queue --url http://localhost:15672/api/queues --username monitoring --password password --vhost Some_Virtual_Host

package main

import (
  "fmt"
  "flag"
  "strconv"
  "os"
)

var _ = fmt.Printf

// Multiple string arguments
type arrayFlags []string

func (i *arrayFlags) String() string {
  return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
  *i = append(*i, value)
  return nil
}

// Multiple int arguments
type arrayFlagsInt []int

func (i *arrayFlagsInt) String() string {
  return "my int representation"
}

func (i *arrayFlagsInt) Set(value string) error {
  value_converted, err := strconv.Atoi(value)
  if err != nil {
    os.Exit(1)
  }

  *i = append(*i, value_converted)
  return nil
}

var QueuesIgnoreFlags arrayFlags
var QueuesFlags arrayFlags

var WarningThresholdFlag arrayFlagsInt
var CriticalThresholdFlag arrayFlagsInt

var URLFlag string
var UsernameFlag string
var PasswordFlag string
var VHostFlag string

var DefaultWarningThresholdFlag int
var DefaultCriticalThresholdFlag int

func main() {
  flag.Var(&QueuesIgnoreFlags, "queues-ignore", "test")
  flag.Var(&QueuesFlags, "queue", "This is a queue name. You can add multiple queues with multiple --queue arguments")

  flag.Var(&WarningThresholdFlag, "warning-threshold", "test")
  flag.Var(&CriticalThresholdFlag, "critical-threshold", "test")

  flag.StringVar(&URLFlag, "url", "", "This is RabbitMQ URL")
  flag.StringVar(&UsernameFlag, "username", "", "This is RabbitMQ API username")
  flag.StringVar(&PasswordFlag, "password", "", "This is RabbitMQ API password")
  flag.StringVar(&VHostFlag, "vhost", "", "This is RabbitMQ virtual host")

  flag.IntVar(&DefaultWarningThresholdFlag, "default-warning-threshold", 0, "test")
  flag.IntVar(&DefaultCriticalThresholdFlag, "default-critical-threshold", 0, "test")

  // Parse all arguments in
  flag.Parse()

  fmt.Println("Queues to be ignored:", QueuesIgnoreFlags)
  fmt.Println("Queues to be monitored:", QueuesFlags)
  fmt.Println("Warning thresholds for these queues:", WarningThresholdFlag)
  fmt.Println("Critical thresholds for these queues:", CriticalThresholdFlag)

  fmt.Println("URL:", URLFlag)
  fmt.Println("Username:", UsernameFlag)
  fmt.Println("Password:", PasswordFlag)
  fmt.Println("Virtual host:", VHostFlag)

  fmt.Println("Default warning threshold:", DefaultWarningThresholdFlag)
  fmt.Println("Default critical threshold:", DefaultCriticalThresholdFlag)
}

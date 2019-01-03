//usr/bin/env go run "$0" "$@"; exit

// An example to call this script (remove loop function)
//
// loop ./monitoring.go --queue-ignore ignore_this_queue --queue-ignore ignore_this_queue_as_well --threshold-warning 1 --threshold-critical 3 --threshold-warning 2 --threshold-critical 5 --threshold-warning-default 4 --threshold-critical-default 5 --queue some_incoming_queue --queue some_outgoing_queue --api-url 'http://localhost:15672/api/queues' --api-username monitoring --api-password password --vhost Some_Virtual_Host --debug no

package main

import (
  "fmt"
  "flag"
  "strconv"
  "os"
  "io/ioutil"
  "net/http"
  "log"
  "encoding/json"
  "time"
)

var _ = fmt.Printf

// Arguments parsing
type arrayFlags []string

func (i *arrayFlags) String() string {
  return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
  // Add to the array
  *i = append(*i, value)
  return nil
}

// Multiple int arguments
type arrayFlagsInt []int

func (i *arrayFlagsInt) String() string {
  return "my int representation"
}

func (i *arrayFlagsInt) Set(value string) error {
  // Convert string to integer
  value_converted, err := strconv.Atoi(value)
  if err != nil {
    os.Exit(1)
  }

  // Add to the array
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
var DebugFlag string

var DefaultWarningThresholdFlag int
var DefaultCriticalThresholdFlag int

// JSON parsing
type PublicKey struct {
  Name string
  Messages int
  Vhost string
}

// Function without return, because all variables are global
func parse_args() {
  // Parsing arguments
  flag.Var(&QueuesIgnoreFlags, "queue-ignore", "You need to specify the queues that should be excluded from monitoring, otherwise they will be monitored, e.g. --queue-ignore queue1 --queue-ignore queue2")
  flag.Var(&QueuesFlags, "queue", "This is a queue name for monitoring. You can add multiple queues, e.g. --queue queueA --queue queueB")

  flag.Var(&WarningThresholdFlag, "threshold-warning", "This is a warning threshold. If you have specified multiple queues, you need multiple --threshold-warning parameters, e.g. --threshold-warning 10 --threshold-warning 20")
  flag.Var(&CriticalThresholdFlag, "threshold-critical", "This is a critical threshold. If you have specified multiple queues, you need multiple --threshold-critical parameters, e.g. --threshold-critical 20 --threshold-critical 20")

  flag.StringVar(&URLFlag, "api-url", "", "This is RabbitMQ URL, e.g. --api-url 'http://10.10.10.10:15672/api/queues'")
  flag.StringVar(&UsernameFlag, "api-username", "", "This is RabbitMQ API username, e.g. --api-username monitoring")
  flag.StringVar(&PasswordFlag, "api-password", "", "This is RabbitMQ API password, e.g. --api-password MyPassword")
  flag.StringVar(&VHostFlag, "vhost", "", "This is your target RabbitMQ virtual host, but IMPORTANT that other virtual hosts will not be monitored")
  flag.StringVar(&DebugFlag, "debug", "no", "For more verbose output use --debug yes")

  flag.IntVar(&DefaultWarningThresholdFlag, "threshold-warning-default", 40, "If you didn't specify the queue it's still monitored and has warning threshold")
  flag.IntVar(&DefaultCriticalThresholdFlag, "threshold-critical-default", 50, "If you didn't specify the queue it's still monitored and has critical threshold")

  // Parse all arguments in
  flag.Parse()

  // Check whether the same count of arguments passed for each queue
  if len(QueuesFlags) != len(WarningThresholdFlag) || len(QueuesFlags) != len(CriticalThresholdFlag) {
    fmt.Println("It seems you have passed different count of arguments for queues and queues thresholds, please check --help")
    os.Exit(3)
  }

  // Debug mode
  if DebugFlag == "yes" {
    fmt.Println("---- Arguments ----")
    fmt.Printf("%v queue(s) to be ignored: %v\n", len(QueuesIgnoreFlags), QueuesIgnoreFlags)
    fmt.Printf("%v queue(s) to be monitored: %v\n", len(QueuesFlags), QueuesFlags)
    fmt.Printf("%v warning threshold(s) for the queue(s): %v\n", len(WarningThresholdFlag), WarningThresholdFlag)
    fmt.Printf("%v critical threshold(s) for the queue(s): %v\n", len(CriticalThresholdFlag), CriticalThresholdFlag)

    fmt.Println("URL:", URLFlag)
    fmt.Println("Username:", UsernameFlag)
    fmt.Println("Password:", PasswordFlag)
    fmt.Println("Virtual host:", VHostFlag)

    fmt.Println("Default warning threshold:", DefaultWarningThresholdFlag)
    fmt.Println("Default critical threshold:", DefaultCriticalThresholdFlag)
    fmt.Println("---- ---- ---- ----")
  }
}

func http_query(method string, address string, user string, pass string) string {
  // Use HTTP connection
  var client_with_timeout = &http.Client{
    Timeout: time.Second * 10,
  }

  // Doing actual HTTP request
  request, error := http.NewRequest(method, address, nil)

  // Basic HTTP authentication
  request.SetBasicAuth(user, pass)

  if error != nil{
    log.Fatal(error)
  }

  // Getting the output
  response, error := client_with_timeout.Do(request)

  if error != nil{
    log.Fatal(error)
  }

  // Checking HTTP response
  defer response.Body.Close()

  if response.StatusCode != 200 {
    log.Fatal("HTTP response code is not 200, for example: incorrect user's credentials, access denied, etc.")
  }

  // Write the output to variable (consumes memory based on JSON size)
  bodyText, error := ioutil.ReadAll(response.Body)
  response.Body.Close()
  if error != nil {
    log.Fatal(error)
  }

  // Getting the output
  var result map[string]interface{}
  json.Unmarshal([]byte(bodyText), &result)

  // We have whole JSON in "[]byte" data type and want to convert to string
  bodyText_string := string(bodyText)
  return bodyText_string
}

// There is no built-in operator in Go to check whether array contains a string so writing my own
func contains(arr []string, str string) bool {
  for _, a := range arr {
    if a == str {
      return true
    }
  }
  return false
}

func parse_json(whole_json string) string {
  // Parse JSON
  input := []byte(whole_json)
  var output []PublicKey
  json.Unmarshal([]byte(input), &output)

  // This is the variable to check whether at least 1 queue found for monitoring
  any_queues := false

  // We need to count queues to track the order
  queue_counter := 0

  // We have thresholds specified 1)by user or 2)default, these variables to combine them into one variable later
  warning_threshold := 0
  critical_threshold := 0

  // Alert enable/disable
  warning_alert := false
  critical_alert := false

  // Go through all queues from a whole JSON
  for _ , from_json := range output {
    // Skip if vhost not matched
    if VHostFlag != from_json.Vhost {
      fmt.Printf("%v (skip this, vhost not matched)\n", from_json.Name)
      continue
    }

    // Skip if queue asked to be ignored
    if contains(QueuesIgnoreFlags, from_json.Name) {
      fmt.Printf("%v (skip this, asked to be ignored)\n", from_json.Name)
      continue
    }

    // At this stage we know that there is at least 1 queue for monitoring
    any_queues = true

    // Match queues in JSON that originally were specified by user
    if contains(QueuesFlags, from_json.Name) {
      // Assign to common *_threshold parameters
      // Get thresholds based on order, passed as the arguments by the user
      warning_threshold = WarningThresholdFlag[queue_counter]
      critical_threshold = CriticalThresholdFlag[queue_counter]
      queue_counter = queue_counter + 1

      if DebugFlag == "yes" {
        fmt.Printf("Monitor: %v\n", from_json.Name)
      }
    } else {
      // Here goes queues that were not specified by user, they get default thresholds
      warning_threshold = DefaultWarningThresholdFlag
      critical_threshold = DefaultCriticalThresholdFlag

      if DebugFlag == "yes" {
        fmt.Printf("Monitor additional: %v\n", from_json.Name)
      }
    }
    // Print all queues for monitoring
    fmt.Printf("%v Current: %v Warning: %v Critical: %v\n", from_json.Name, from_json.Messages, warning_threshold, critical_threshold)

    // Collect alert information
    if from_json.Messages >= critical_threshold {
      critical_alert = true
    } else if from_json.Messages >= warning_threshold {
      warning_alert = true
    }
  }

  // ---- Exit the script for critical alert ----
  if critical_alert {
    fmt.Println("! critical alert !")
    os.Exit(101)
  }

  // ---- Exit the script for warning alert, only if there was no critical alert ----
  if warning_alert {
    fmt.Println("! warning alert !")
    os.Exit(100)
  }

  // Tell the user when nothing is found for monitoring
  if any_queues == false {
    fmt.Printf("No queues found for virtual host named %v on %v\n", VHostFlag, URLFlag)
  }

  // If script ran to this line, it means there are no alerts at this moment
  return "exit status 0"
}

func main() {
  // Parse all arguments (they will be asigned as global variables)
  parse_args()

  // Do HTTP query
  http_query_content := http_query("GET", URLFlag, UsernameFlag, PasswordFlag)

  // Parse JSON and throw alert if threshold reached
  am_i_successful := parse_json(http_query_content)

  // Just an output from the function
  fmt.Println(am_i_successful)
}

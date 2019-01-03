//usr/bin/env go run "$0" "$@"; exit

// An example to call this script (remove loop function)
//
// loop ./monitoring.go --queues-ignore ignore_this_queue --queues-ignore ignore_this_queue_as_well --warning-threshold 1 --critical-threshold 3 --warning-threshold 2 --critical-threshold 4 --default-warning-threshold 4 --default-critical-threshold 5 --queue some_incoming_queue --queue some_outgoing_queue --url http://localhost:15672/api/queues --username monitoring --password password --vhost Some_Virtual_Host

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
  flag.Var(&QueuesIgnoreFlags, "queues-ignore", "test")
  flag.Var(&QueuesFlags, "queue", "This is a queue name. You can add multiple queues with multiple --queue arguments")

  flag.Var(&WarningThresholdFlag, "warning-threshold", "test")
  flag.Var(&CriticalThresholdFlag, "critical-threshold", "test")

  flag.StringVar(&URLFlag, "url", "", "This is RabbitMQ URL")
  flag.StringVar(&UsernameFlag, "username", "", "This is RabbitMQ API username")
  flag.StringVar(&PasswordFlag, "password", "", "This is RabbitMQ API password")
  flag.StringVar(&VHostFlag, "vhost", "", "This is RabbitMQ virtual host")
  flag.StringVar(&DebugFlag, "debug", "no", "Set it to 'yes' for more verbose output")

  flag.IntVar(&DefaultWarningThresholdFlag, "default-warning-threshold", 0, "test")
  flag.IntVar(&DefaultCriticalThresholdFlag, "default-critical-threshold", 0, "test")

  // Parse all arguments in
  flag.Parse()

  // Debug mode
  if DebugFlag == "yes" {
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

// There is no built-in operator to check whether array contains a string so writing my own
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

  // We have thresholds specified by user or default, these variables to combine them
  warning_threshold := 0
  critical_threshold := 0

  // Alert enable/disable
  warning_alert := false
  critical_alert := false

  // Go through all queues in a whole JSON
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

    // At this stage we know that we have at least 1 queue for monitoring
    any_queues = true

    // Match queues in JSON that originally were specified by user
    if contains(QueuesFlags, from_json.Name) {
      // Assign to common *_threshold parameters, also get thresholds based on order originally specified by user
      // fmt.Printf("Monitor: %v\n", from_json.Name)
      warning_threshold = WarningThresholdFlag[queue_counter]
      critical_threshold = CriticalThresholdFlag[queue_counter]
      queue_counter = queue_counter + 1
    } else {
      // Here goes queues that were not specified by user - "default thresholds"
      // fmt.Printf("Monitor additional: %v\n", from_json.Name)
      warning_threshold = DefaultWarningThresholdFlag
      critical_threshold = DefaultCriticalThresholdFlag
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
    fmt.Printf("! critical alert !\n")
    os.Exit(101)
  }

  // ---- Exit the script for warning alert, the priority is for critical over warning alert ----
  if warning_alert {
    fmt.Printf("! warning alert !\n")
    os.Exit(100)
  }

  // Tell the user when nothing is found for monitoring
  if any_queues == false {
    fmt.Printf("No queues found for virtual host named %v on %v\n", VHostFlag, URLFlag)
  }

  // If script went to this phase, it means there are no alerts at this moment
  return "exit status 0"
}

  func main() {
    // Parse all arguments (they will be asigned as global variables)
    parse_args()

    // Do HTTP query
    http_query_content := http_query("GET", URLFlag, UsernameFlag, PasswordFlag)

    // Parse JSON and get required elements
    am_i_successful := parse_json(http_query_content)

    fmt.Println(am_i_successful)
  }

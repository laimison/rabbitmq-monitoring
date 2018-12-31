//usr/bin/env go run "$0" "$@"; exit

// An example to call this script (remove loop function)
//
// loop ./monitoring.go --queues-ignore ignore_this_queue --queues-ignore ignore_this_queue_as_well --warning-threshold 1 --critical-threshold 2 --warning-threshold 1 --critical-threshold 2 --default-warning-threshold 4 --default-critical-threshold 5 --queue some_incoming_queue --queue some_outgoing_queue --url http://localhost:15672/api/queues --username monitoring --password password --vhost Some_Virtual_Host

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
  // "encoding/base64"
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

var DefaultWarningThresholdFlag int
var DefaultCriticalThresholdFlag int

// JSON parsing
type PublicKey struct {
  Name string
  Messages int
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

  flag.IntVar(&DefaultWarningThresholdFlag, "default-warning-threshold", 0, "test")
  flag.IntVar(&DefaultCriticalThresholdFlag, "default-critical-threshold", 0, "test")

  // Parse all arguments in
  flag.Parse()

  // Output to the screen
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

func http_query(method string, address string, user string, pass string) string {
  // Use HTTP connection
  var client_with_timeout = &http.Client{
    Timeout: time.Second * 10,
  }

  // Doing actual HTTP request
  request, error := http.NewRequest(method, address, nil)

  if error != nil{
    log.Fatal(error)
  }

  // Basic HTTP authentication
  request.SetBasicAuth(user, pass)
  // user, pass, ok := request.BasicAuth()
  // fmt.Println(ok)

  // Getting the output
  response, error := client_with_timeout.Do(request)

  if error != nil{
    log.Fatal(error)
  }

  // Check if password was correct
  // auth := response.Header.Get("Authorization")
  // up, _ := base64.StdEncoding.DecodeString(auth[6:])
  //
  // fmt.Println(up)

  // Write the output to variable (memory)
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

func main() {
  // Parse all arguments
  parse_args()

  // Do HTTP query
  http_query_content := http_query("GET", URLFlag, UsernameFlag, "password")

  input := []byte(http_query_content)

  var output []PublicKey
  json.Unmarshal([]byte(input), &output)

  for _ , value := range output {
    fmt.Printf("%v: %v\n", value.Name, value.Messages)
  }
}

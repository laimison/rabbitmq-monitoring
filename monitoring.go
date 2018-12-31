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
)

var _ = fmt.Printf

// Multiple string arguments
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

func http_query(method string, address string, user string, pass string) ([]byte) {
  // Use HTTP connection
  var client_with_timeout = &http.Client{
    Timeout: time.Second * 10,
  }

  // Doing actual HTTP request
  request, error := http.NewRequest(method, address, nil)

  // Basic HTTP authentication
  request.SetBasicAuth(user, pass)

  // Getting the output
  response, error := client_with_timeout.Do(request)

  // Throw an error
  if error != nil{
    log.Fatal(error)
  }

  // Write the output to memory - variable
  bodyText, error := ioutil.ReadAll(response.Body)

  // Getting the output
  var result map[string]interface{}
  json.Unmarshal([]byte(bodyText), &result)

  // Output whole JSON in "[]byte" data type
  return bodyText
}

func parse_json_example() {
  // Parsing json
  type Bird struct {
    Species string
    Description string
  }

  birdJson := `{"species": "pigeon","description": "likes to perch on rocks"}`
  var bird Bird

  json.Unmarshal([]byte(birdJson), &bird)

  fmt.Printf("\nSpecies:\n%s\n\nDescription:\n%s\n", bird.Species, bird.Description)
}

func parse_json_example2() {
  jsonString := http_query("GET", "https://api.github.com/users/1", "", "")

  var jsonMap map[string]interface{}
  json.Unmarshal([]byte(jsonString ), &jsonMap)

  fmt.Println(jsonMap)
}

func main() {
  parse_args()

  // Doing HTTP query
  http_query_content := http_query("GET", "http://localhost:15672/api/queues", "monitoring", "password")
  // fmt.Println(string(http_query_content))

  type PublicKey2 struct {
    Name string
    Messages int
  }

  input2 := []byte(string(http_query_content))

  var output2 []PublicKey2
  json.Unmarshal([]byte(input2), &output2)

  for key2, value2 := range output2 {
    fmt.Printf("Id: %v\n", key2)
    fmt.Printf("Name: %v\n", value2.Name)
    fmt.Printf("Messages: %v\n", value2.Messages)
  }

  type PublicKey struct {
    Id int
    People string
    // [] tells that it's json array, remove if you don't have array in front of groups:
    Groups []struct {
      Name string
    }
  }

  input := []byte(`[
    {"id": 1, "people": "Tom", "groups": [{"name": "groupA"}]},
    {"id": 2, "people": "Ben", "groups": [{"name": "groupB"}]},
    {"id": 3, "people": "Dan", "groups": [{"name": "groupA"}, {"name": "groupC"}]}
    ]`)

  var output []PublicKey
  json.Unmarshal([]byte(input), &output)

  for key, value := range output {
    fmt.Printf("Id: %v\n", key)
    fmt.Printf("People: %v\n", value.People)

    // fmt.Printf("Groups: %v\n", value.Groups[0].Name)
    for key_groups, value_groups := range value.Groups {
      fmt.Printf("  Group %v: %v\n", key_groups, value_groups.Name)
    }
    fmt.Println()
  }

  // b, err := json.MarshalIndent(http_query_content, "", "   ")
  // if err != nil {
  //   fmt.Println("error:", err)
  // }
  // os.Stdout.Write(b)

  // parse_json_example()
  // parse_json_example2()
}

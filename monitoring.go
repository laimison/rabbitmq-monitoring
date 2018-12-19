//usr/bin/env go run "$0" "$@"; exit

package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
    "encoding/json"
    "time"
    // "bytes"
    // "os"
)

func rest_query(method string, address string, user string, pass string) string {
  // HTTP client to connect to API
  // client := &http.Client{}
  var client_timeout = &http.Client{
    Timeout: time.Second * 10,
  }

  // Doing query
  request, error := http.NewRequest(method, address, nil)

  // Basic HTTP authentication
  request.SetBasicAuth(user, pass)

  // Getting the output
  response, error := client_timeout.Do(request)

  // Throw error
  if error != nil{
      log.Fatal(error)
  }

  // Write the output to memory
  bodyText, error := ioutil.ReadAll(response.Body)

  // Getting the output
  var result map[string]interface{}
  json.Unmarshal([]byte(bodyText), &result)

  // fmt.Println(result)
  content := string(bodyText)
  // content := string(result)
  return content
}

func main() {
  // Doing HTTP query
  rest_query_content := rest_query("GET", "http://localhost:15672/api/queues", "monitoring", "password")
  fmt.Println(rest_query_content)

  // Parsing json
  type Bird struct {
    Species string
    Description string
  }

  birdJson := `{"species": "pigeon","description": "likes to perch on rocks"}`
  var bird Bird

  json.Unmarshal([]byte(birdJson), &bird)

  fmt.Printf("\nSpecies:\n%s\n\nDescription:\n%s\n", bird.Species, bird.Description)

  // Converting string to json
  jsonString := "{\"foo\":{\"baz\": [1,2,3]}}"

  var jsonMap map[string]interface{}
  json.Unmarshal([]byte(jsonString ), &jsonMap)

  fmt.Println(jsonMap) 
}

//usr/bin/env go run "$0" "$@"; exit

package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "log"
    // "bytes"
    // "encoding/json"
    // "os"
)

func rest_query(method string, address string, user string, pass string) string {
  client := &http.Client{}

  request, error := http.NewRequest(method, address, nil)

  request.SetBasicAuth(user, pass)

  response, error := client.Do(request)

  if error != nil{
      log.Fatal(error)
  }

  bodyText, error := ioutil.ReadAll(response.Body)
  content := string(bodyText)

  return content
}

func main() {
  content := rest_query("GET", "http://localhost:15672/api/queues", "monitoring", "password")
  fmt.Println(content)
}

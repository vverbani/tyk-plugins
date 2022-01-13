package main

import (
  "io/ioutil"
  "net/http"
  "log"
  "fmt"

  "github.com/TykTechnologies/tyk/ctx"
)

func mylogger(s string) {
  prefix := "##################################### "
  log.Println(prefix + s)
}

// AddHeader adds custom "Foo: Bar" header to the request
func TestBlocking(rw http.ResponseWriter, r *http.Request) {
  apiDefinition := ctx.GetDefinition(r)
  mylogger("Start")
  for k, v := range apiDefinition.ConfigData {
    if k == "URL" {
      // connect to the url and get the data back
      response, err := http.Get(fmt.Sprintf("%s", v))
      if err != nil {
        mylogger(err.Error())
      }
      responseData, err := ioutil.ReadAll(response.Body)
      if err != nil {
        mylogger(err.Error())
      }
      mylogger(string(responseData))
    }
    mylogger("End")
  }
}

func main() {}

package main

import (
  "net/http"
  "log"
  "fmt"
  "github.com/davecgh/go-spew/spew"

  "github.com/TykTechnologies/tyk/ctx"
)

func mylogger(s string) {
  prefix := "##################################### "
  log.Println(prefix + s)
}

// AddHeader adds custom "Foo: Bar" header to the request
func AddHeadersFromConfigData(rw http.ResponseWriter, r *http.Request) {
  apiDefinition := ctx.GetDefinition(r)
  mylogger("")
  log.Println("%v", apiDefinition.ConfigData)
  log.Print(spew.Sdump(apiDefinition.ConfigData))
  r.Header.Add("Fo", "Bar-mitzvah")
  for k, v := range apiDefinition.ConfigData {
    r.Header.Add(k, fmt.Sprintf("%s", v))
  }
  r.Header.Add("Foo", "Bar-mitzvah")
}

func main() {}

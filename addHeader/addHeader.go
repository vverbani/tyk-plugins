package main

import (
  "net/http"
  "os"
  "bufio"
  "strings"
  "log"
)

func mylogger(s string) {
  prefix := "##################################### "
  log.Println(prefix + s)
}


// AddHeader adds custom "Foo: Bar" header to the request
func AddHeader(rw http.ResponseWriter, r *http.Request) {
  r.Header.Add("Foo", "Bar-mitzvah")
}

func main() {}

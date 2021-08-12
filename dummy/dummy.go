package main

import (
  "net/http"
  "log"
)

func mylogger(s string) {
  prefix := "##################################### "
  log.Println(prefix + s)
}

// DummyOne function just logs one line and exits
func DummyOne(rw http.ResponseWriter, r *http.Request) {
  mylogger("DummyOne Called")
}

// DummyTwo function just logs one line and exits
func DummyTwo(rw http.ResponseWriter, r *http.Request) {
  mylogger("DummyTwo Called")
}

func main() {}

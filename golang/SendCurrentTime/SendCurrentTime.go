package main

import (
    "encoding/json"
    "net/http"
    "time"
    )

func SendCurrentTime(rw http.ResponseWriter, r *http.Request) {
  // check if we don't need to send reply
  if r.URL.Query().Get("get_time") != "1" {
    // allow request to be processed and sent to upstream
    r.Header.Add("Foo", "Bar")
    return
  }

  // prepare data to send
  replyData := map[string]interface{}{
     "current_time": time.Now(),
  }

  jsonData, err := json.Marshal(replyData)
  if err != nil {
    rw.WriteHeader(http.StatusInternalServerError)
    return
  }

  // send HTTP response from Golang plugin
  rw.Header().Set("Content-Type", "application/json")
  rw.WriteHeader(http.StatusOK)
  rw.Write(jsonData)
}

func main() {}

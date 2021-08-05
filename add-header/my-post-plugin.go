package main

import (
  "net/http"
  "os"
  "bufio"
  "strings"
  "log"
)

// AddHeader adds custom "Foo: Bar" header to the request
func AddHeader(rw http.ResponseWriter, r *http.Request) {
  file, err := os.Open("/opt/tyk-plugins/my-post-plugin.txt")
  if err != nil {
    // warning. This causes the gateway to exit
    log.Fatalf("failed opening file: %s", err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

  for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
  }

  for _, eachline := range txtlines {
    text := strings.Split(eachline, ": ")
    r.Header.Add(text[0], text[1])
		//fmt.Println(eachline)
  }

  //r.Header.Add("Foo", "Bar-mitsva")
}

func main() {}

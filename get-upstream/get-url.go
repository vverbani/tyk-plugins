package main

import (
	"log"
	"net/http"
)

// ReportURL logs the upstream URL after it has been rewritten
func ReportURL(rw http.ResponseWriter, r *http.Request) {
	log.Println("Upstream URL is: ", r.URL)
	//r.Header.Add("Foo", "Bar-mitsva")
}

func main() {}

package main

import (
	"log"
	"net/http"
  "github.com/TykTechnologies/tyk/ctx"
  "github.com/davecgh/go-spew/spew"
)

// DumpAPIdefn dumps out the API def visible to the pulgin so we can see what we've got
func DumpAPIdefn(rw http.ResponseWriter, r *http.Request) {

  apidef := ctx.GetDefinition(r)
  log.Println("Start API definition dump")
  spew.Dump(apidef)
  log.Println("End API definition dump")
  if v := r.Context().Value(ctx.UrlRewriteTarget) ; v != nil {
    log.Println(v)
  } else {
    log.Println("FAILED")
    }

}

func main() {}

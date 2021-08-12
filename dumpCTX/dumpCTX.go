package main

import (
	"log"
	"net/http"
  "reflect"

	"github.com/TykTechnologies/tyk/ctx"
	"github.com/davecgh/go-spew/spew"
)

func spewVal(c ctx.Key, s string, r *http.Request) {
	var v string
	if v := r.Context().Value(c); v == nil {
		v = "No value defined"
	}
	log.Println("Start ", s)
	log.Print(spew.Sdump(v))
	log.Println("End ", s)
	return
}

func mylogger(s string) {
  prefix := "##################################### "
  log.Println(prefix + s)
}

// DumpCTX dumps out context variables avialble to the plugin
func DumpCTX(rw http.ResponseWriter, r *http.Request) {

  // Dump the whole http.Request
	mylogger("Start http.Request dump")
	log.Print(spew.Sdump(r))
	mylogger("End http.Request dump")


	// API definition object
	mylogger("Start API definition dump")
	apidef := ctx.GetDefinition(r)
	log.Print(spew.Sdump(apidef))
	mylogger("End API definition dump")

	// Auth Token
	mylogger("Start Auth Token dump")
	authToken := ctx.GetAuthToken(r)
	log.Print(spew.Sdump(authToken))
	mylogger("End Auth Token dump")

	// SessionState
	mylogger("Start Session State dump")
	sessionState := ctx.GetSession(r)
	log.Print(spew.Sdump(sessionState))
	mylogger("End Session State dump")
	mylogger("Start Session Alias")
  log.Print(sessionState.Alias)
	mylogger("End Session Alias")

	// work through the rest of the constants and dump them
	spewVal(ctx.UpdateSession, "ctx.UpdateSession", r)
	spewVal(ctx.HashedAuthToken, "ctx.HashedAuthToken", r)
	spewVal(ctx.VersionData, "ctx.VersionData", r)
	spewVal(ctx.VersionDefault, "ctx.VersionDefault", r)
	spewVal(ctx.OrgSessionContext, "ctx.OrgSessionContext", r)
	spewVal(ctx.ContextData, "ctx.ContextData", r)
	spewVal(ctx.RetainHost, "ctx.RetainHost", r)
	spewVal(ctx.TrackThisEndpoint, "ctx.TrackThisEndpoint", r)
	spewVal(ctx.DoNotTrackThisEndpoint, "ctx.DoNotTrackThisEndpoint", r)
	spewVal(ctx.UrlRewritePath, "ctx.UrlRewritePath", r)
	spewVal(ctx.RequestMethod, "ctx.RequestMethod", r)
	spewVal(ctx.OrigRequestURL, "ctx.OrigRequestURL", r)
	spewVal(ctx.LoopLevel, "ctx.LoopLevel", r)
	spewVal(ctx.LoopLevelLimit, "ctx.LoopLevelLimit", r)
	spewVal(ctx.ThrottleLevel, "ctx.ThrottleLevel", r)
	spewVal(ctx.ThrottleLevelLimit, "ctx.ThrottleLevelLimit", r)
	spewVal(ctx.Trace, "ctx.Trace", r)
	spewVal(ctx.CheckLoopLimits, "ctx.CheckLoopLimits", r)
	spewVal(ctx.UrlRewriteTarget, "ctx.UrlRewriteTarget", r)
	spewVal(ctx.TransformedRequestMethod, "ctx.TransformedRequestMethod", r)
	spewVal(ctx.RequestStatus, "ctx.RequestStatus", r)
	// Only in 3.0+
	//spewVal(ctx.GraphQLRequest, "ctx.GraphQLRequest", r)
	mylogger("Start Metadata")
	log.Println(reflect.ValueOf(sessionState).Elem().FieldByName("MetaData"))
	for key, value := range sessionState.MetaData {
		log.Print(key, "->", value)
	}
	mylogger("End Metadata")
}

func main() {}

package main

import (
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/zhaostu/grj/transit_realtime"
)

var port = flag.Int("port", 8000, "The port the server listens on.")

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	q, ok := query["q"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uri := q[0]
	if uri == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the body for GTFS realtime feed
	resp, err := http.Get(uri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unpack the protobuf to message
	var msg transit_realtime.FeedMessage
	err = proto.Unmarshal(data, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pack the message to json
	w.Header().Set("Content-Type", "application/json")
	marshaler := jsonpb.Marshaler{Indent: "  "}
	marshaler.Marshal(w, &msg)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

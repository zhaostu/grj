package grj

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"grj/transit_realtime"
)

func init() {
	http.HandleFunc("/", handler)
}

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
	// App engine stuff
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	resp, err := client.Get(uri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Unpack the protobuf to message
	var msg transit_realtime.FeedMessage
	err = proto.Unmarshal(data, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pack the message to json
	w.Header().Set("Content-Type", "application/json")
	marshaler := jsonpb.Marshaler{Indent: "  "}
	marshaler.Marshal(w, &msg)
}

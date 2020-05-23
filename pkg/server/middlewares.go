package server

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"io/ioutil"
	"monkiato/apio-orchestrator/internal/db"
	"net/http"
)

// ParseBody middleware used to apply a JSON parse for the request body. Data will be stored in Gorilla Context
// it can be obtained from subsequence handlers through context.Get(r, "parseBody")
func ParseBody(handler http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			addErrorResponse(w, http.StatusBadRequest, "can't read body")
			return
		}
		var parsedBody map[string]interface{}
		jsonErr := json.Unmarshal(body, &parsedBody)
		if jsonErr != nil {
			addErrorResponse(w, http.StatusBadRequest, "can't parse body")
			return
		}
		context.Set(r, "parsedBody", parsedBody)
		handler.ServeHTTP(w, r)
	}
}

// ValidateID middleware used to detect an node ID in the request, if exists it means the endpoint is trying to operate
// over an existing node, and the middleware will try to find and get the node from the storage/DB, otherwise an error
// is returned if the node was not found. The node data will be stored in Gorilla Context, it can be obtained from
// subsequence handlers through context.Get(r, "node")
func ValidateID() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			if nodeId, ok := vars["id"]; ok {
				// id exists, validate if the node exists in the DB
				node, found := db.GetNode(nodeId)
				if !found {
					//not found
					addErrorResponse(w, http.StatusNotFound, "node not found")
					return
				}

				context.Set(r, "node_id", nodeId)
				context.Set(r, "node", node)
			}
			next.ServeHTTP(w, r)
		})
	}
}

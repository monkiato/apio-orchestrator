package server

import (
	"encoding/json"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateNodeHandler(w http.ResponseWriter, r *http.Request) {
	// handle PUT request to create a new node
	node := context.Get(r, "node")
	//TODO: parse data using node.Metadata and save in DB

	//TODO: deploy docker container
	data, err := json.Marshal(node)
	if err != nil {
		log.Error(err.Error())
		addErrorResponse(w, http.StatusInternalServerError, "unable to parse node data")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func EditNodeHandler(w http.ResponseWriter, r *http.Request) {
	// handle POST request to update node data (also for adding collections to the node)
	item := context.Get(r, "item")
	//itemId := context.Get(r, "item_id")
	//TODO: update existing node using the info already provided by the middleware
	//TODO: rename item to node

	//TODO: redeploy docker container associated to this node
	data, err := json.Marshal(item)
	if err != nil {
		log.Error(err.Error())
		addErrorResponse(w, http.StatusInternalServerError, "unable to parse node data")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func RemoveNodeHandler(w http.ResponseWriter, r *http.Request) {
	// handle DELETE request to remove node
	nodeId := context.Get(r, "node_id")
	//TODO: remove existing node from DB

	//TODO: remove node from docker
	//data, err := json.Marshal(item)
	//if err != nil {
	//	log.Error(err.Error())
	//	addErrorResponse(w, http.StatusInternalServerError, "unable to parse node data")
	//	return
	//}
	//w.WriteHeader(http.StatusOK)
	//w.Write(data)
}
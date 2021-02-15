package main

import (
	"fmt"
	"net/http"
	"github.com/kirannhegde/Container-Orchestrator/db"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	initDBConn()
	initClusterNodes()
	setupRoutes(router)
	http.ListenAndServe(":8080", router)
}

func setupRoutes(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/create_container/{registry:[\\w/\\.]+}/{repository}/{image}/{replicas}", createContainerHandler).Methods("POST")
	r.HandleFunc("/state", getClusterState).Methods("GET")
}

func initDBConn() {
	var err error
	db.DBConn, err = gorm.Open("sqlite3", "containers.db")
	if err != nil {
		panic("Failed to connect to the containers db")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page")
}

func createContainerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating containers...")
	vars := mux.Vars(r)
	containerRegistry := vars["registry"]
	containerRepository := vars["repository"]
	containerImage := vars["image"]
	numOfContainerReplicas := vars["replicas"]
	fmt.Fprintf(w, "Registry:%v, repository:%v, image:%v, replicas:%v", containerRegistry, containerRepository, containerImage, numOfContainerReplicas)

	/* urlParams := r.URL.Query()
	for key, value := range urlParams {
		fmt.Fprintf(w, "Key=%v, value=%v", key, value)
	} */

}

func getClusterState(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting cluster state...")
}

func initClusterNodes() {
	fmt.Printf("Initializing the cluste nodes...")
}

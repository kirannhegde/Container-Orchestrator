package main

import (
	"fmt"
	"net/http"
	"github.com/kirannhegde/Container-Orchestrator/db"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

//Entry point to the container orchestrator
func main() {
	router := mux.NewRouter().StrictSlash(true)
	initDBConn()
	initClusterNodes()
	setupRoutes(router)
	http.ListenAndServe(":8080", router)
}

//RESTAPI routes should be defined here.
//The Gorilla mux router from the Gorilla web toolkit  is used here, as it's a more sophisticated router
//than the default router provided by the net/http package.
func setupRoutes(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/create_container/{registry:[\\w/\\.]+}/{repository}/{image}/{replicas}", createContainerHandler).Methods("POST")
	r.HandleFunc("/state", getClusterStateHandler).Methods("GET")
}

func initDBConn() {
	//Using SQLLite since it can be embedded into the application(if required) and can be accessed quickly
	var err error
	db.DBConn, err = gorm.Open("sqlite3", "containers.db")
	if err != nil {
		panic("Failed to connect to the containers db")
	}
}

//Handler function for the route:"/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page")
}

//Handler function to handle the creation of new containers.
//Route:"/create_container/<container_registry>/<image name>:[tag]/<number of replicas>
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

//Handler function to return the cluster state.
//Route:"/state"
func getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting cluster state...")
}

func initClusterNodes() {
	fmt.Printf("Initializing the cluster nodes...")
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/kirannhegde/Container-Orchestrator/db"
	"github.com/kirannhegde/Container-Orchestrator/node"
	"github.com/kirannhegde/Container-Orchestrator/container"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

)

//Entry point to the container orchestrator
func main() {
	router := mux.NewRouter().StrictSlash(true)
	//initDBConn()
	//defer db.DBConn.Close()
	initClusterNodes()
	setupRoutes(router)
	http.ListenAndServe(":8080", router)
}

//RESTAPI routes should be defined here.
//The Gorilla mux router from the Gorilla web toolkit  is used here, as it's a more sophisticated router
//than the default router provided by the net/http package.
func setupRoutes(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
	//Ideal rest end point. Read the comments against the function:createContainerHandler
	//r.HandleFunc("/create_container/{registry:[\\w/\\.]+}/{repository}/{image}/{replicas}", createContainerHandler).Methods("POST")
	r.HandleFunc("/create_container/{image}/{replicas}", createContainerHandler).Methods("POST")
	r.HandleFunc("/state", getClusterStateHandler).Methods("GET")
}

func initDBConn() {
	//Using SQLLite since it can be embedded into the application(if required) and can be accessed quickly
	var err error
	db.DBConn, err = gorm.Open("sqlite3", "orchestrator.db")
	if err != nil {
		panic("Failed to connect to the containers db")
	}
	//Create the "node" table
	db.DBConn.AutoMigrate(&node.ClusterNode{})
	//Create the "containers" table
	db.DBConn.AutoMigrate(&container.Container{})
}

//Handler function for the route:"/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page")
}

/*Handler function to handle the creation of new containers.
Below mentioned is the ideal REST endpoint. 
Route:"/create_container/<container_registry>/<image name>:[tag]/<number of replicas> 
For v1 of this REST server, i will make an assumption that we are only going to be
working of the images on the Docker public repository. Hence, the endpoint changes to something
like:/create_container/<image name>:[tag]/<number of replicas> */
func createContainerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating containers...")
	vars := mux.Vars(r)
	//containerRegistry := vars["registry"]
	//containerRepository := vars["repository"]
	containerImage := vars["image"]
	numOfContainerReplicas, _ := strconv.Atoi(vars["replicas"])
	container := container.Container{
		RequiredNumOfReplicas: numOfContainerReplicas,
	}
	container.Image.ImageNameTag = containerImage

	err := container.CreateContainers()
	if err != nil {
		panic("There was an issue creating containers.")
	}
	

}

//Handler function to return the cluster state.
//Route:"/state"
func getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting cluster state...")
}

func initClusterNodes() {
	fmt.Printf("Initializing the cluster nodes...")
}

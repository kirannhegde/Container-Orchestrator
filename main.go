package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/kirannhegde/Container-Orchestrator/container"
	"github.com/kirannhegde/Container-Orchestrator/db"
	"github.com/kirannhegde/Container-Orchestrator/node"
	_ "github.com/mattn/go-sqlite3"
)

//Entry point to the application
func main() {
	//initDBConn()
	//defer db.DBConn.Close()
	initClusterNodes()
	setupRoutes()
}

//RESTAPI routes should be defined here.
//The Gorilla mux router from the Gorilla web toolkit  is used here, as it's a more sophisticated router
//than the default router provided by the net/http package.
func setupRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/create", createContainerHandler).Methods("POST")
	router.HandleFunc("/state", getClusterStateHandler).Methods("GET")
	router.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", router)
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

/*Handler function to handle the creation of new containers.
The ideal REST end point for container creation should support the following:
a)container registry
b)image name and tag
c)Number of replicase
For v1 of this REST server, i will make an assumption that we are only going to be
working of the images on the Docker public repository. */
func createContainerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating containers...")
	/*vars := mux.Vars(r)
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
	}*/

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))

}

//Handler function to return the cluster state.
//Route:"/state"
func getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting cluster state...")
}

//Handler function for the route:"/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page")
}

func initClusterNodes() {
	fmt.Printf("Initializing the cluster nodes...")
}

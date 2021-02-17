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

/*Entry point to the container orchestrator
Assumptions:
1)This container orchestrator will only be capable of handling Docker container.
2)The cluster will already have some nodes added to begin with.

So we will be using a SQLLite DB to store the information about various nodes of the cluster,
information about the various containers running on the nodes, node capacity, available capacity
The structure of the database could be something like:
1)Node.db  - This table will have the following fields:
	id           int
	nodeName     string
	nodeIpaddr   string
	nodeCapacity int64	//Represents the total capacity of the node
	nodeAvailCapacity int64	//Represents the available capacity of the node.
	status 		 string //Indicates whether the node is online or offline.

To make things simpler, lets assume that all nodes have the same fixed capacity.
2)container.db - This table will hold data about all the running containers in the cluster.
This table will have the following fields:
	Image //Image from which the container has been created.
	RequiredNumOfReplicas int
	actualNumOfReplicas   int
	nodeName              string
	nodeIpaddr            string

I foresee the need to have multiple go routines running during the life cycle of the RESTAPI server.

***For the time being lets assume that, we have few nodes which are part of the cluster. ***
Assuming that we already have nodes in the cluster,
when the RESTAPI server is started, the following steps should be executed:
1)The connection to the SQLLite DB should be established.
The DB connection is shared by everyone in the server(DBConn *gorm.DB, defined in db.go)
2)The information about the nodes that form the cluster is fetched from the DB(node.db)
3)A ***go routine*** will get spawned whose only function is to check the health of the nodes
and update the status in the DB. We can call this go routine: SendHeartBeatToNodes()
This go routine will be running in the ***background*** throughout the lifecycle of the API server.
The other purpose of this go routine is to update the health of every node in the database on a
periodic basis.
This go routine will also communicate with another go routine(described below), to notify about
nodes going down and hence the need to rebalance containers.  This communication will happen by exchanging
messages on a dedicated "unbuffered channel."
We specifically want to use an "unbuffered channel" so as to ensure that the rebalacing is reliable.
4)We will have another go routine listening for incoming messages
from the SendHeartBeatToNodes() go routine. Lets call this go-routine : RebalanceContainersOnNodes()
Once a message is received from the go routine started in step #3, this go routine will inspect the
node.db table to get the list of nodes which are offline.  With this information in hand, this go routine will then
inspect the containers.db table to rebalance the affected containers.



Briefly the working of the REST API end points will be as follows:
a)/create  POST
Once there is a POST request to create a container with n number of replicas, the following workflow
would kick in:
1)The handler function(createContainerHandler) would get trigerred.
2)The handler would invoke the container.CreateContainers() method in container.go
In this method, we will use Mutexes from the sync.Mutex package or we will use a buffered channel with a capacity of 1
to establish locks before the containers are spawned.
This is to ensure that the spawning of containers is synchronized and container instances of only one image are being spawned
at a time.
3)Once the containers are spawned, the container.db table is updated.

Spawning of containers:
All the containers can be spawned using the Docker remote API.

Additional TODO:
-To interact with containers in container.go, make use of interfaces or rather program against interfaces.
Using interfaces helps you extend your code for other container types easily and also makes the maintenance of the code easier.
-Use interfaces to interact with sqllite db in db.go

*/
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

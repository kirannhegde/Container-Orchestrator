# Container Orchestrator
A simple project in GoLang to mimic a container orchestrator

### High level design

**Assumptions:** <br />
1)This container orchestrator will only be capable of handling Docker container.<br />
2)The cluster will already have some nodes added to begin with.<br />

We will be using a SQLLite DB to store the information about various nodes of the cluster,
information about the various containers running on the nodes, node capacity, available capacity
The structure of the database could be something like:<br />
1)**Node.db** <br /> This table will have the following fields:<br />
	    id  -->         int <br />
    	nodeName    -->     string<br />
    	nodeIpaddr -->       string<br />
    	nodeCapacity -->     int64	//Represents the total capacity of the node<br />
    	nodeAvailCapacity -->     int64	//Represents the available capacity of the node.<br />
    	status 	-->     string //Indicates whether the node is online or offline.<br />
    
To make things simpler, lets assume that all nodes have the same fixed capacity.<br />
2)**container.db**<br /> This table will hold data about all the running containers in the cluster.<br />
This table will have the following fields:<br />
	Image //Image from which the container has been created.<br />
	RequiredNumOfReplicas -->     int<br />
	actualNumOfReplicas  -->      int<br />
	nodeName  -->                string<br />
	nodeIpaddr    -->            string<br />
	
**Important**: To ensure the consistency of the data in the database, we use a single database connection throughout the lifetime<br />
of the RESTAPI server. This same connection is used from various go routines. SQLLite uses a global lock before writing to the db<br />

I foresee the need to have multiple go routines running during the life cycle of the RESTAPI server.<br />

***For the time being lets assume that, we have few nodes which are already part of the cluster.(i.e a static list of nodes)***
Assuming that we already have nodes in the cluster,when the RESTAPI server is started, the following steps should be executed:<br />
1)The connection to the SQLLite DB should be established.<br />
The DB connection is shared by everyone in the server(**DBConn *gorm.DB, defined in db.go**)<br />
2)The information about the nodes that form the cluster is fetched from the DB(node.db)<br />
3)A ***go routine*** will get spawned whose only function is to check the health of the nodes<br />
and update the status in the DB. We can call this go routine: **SendHeartBeatToNodes()**<br />
This go routine will be running in the ***background*** throughout the lifecycle of the API server.<br />
The other purpose of this go routine is to update the health of every node in the database on a periodic basis.<br />
This go routine will also communicate with another go routine(described below), to notify about <br />
nodes going down and hence the need to rebalance containers.  This communication will happen by exchanging<br />
messages on a dedicated "**unbuffered channel.**"<br />
We specifically want to use an "unbuffered channel" so as to ensure that the rebalacing is reliable. <br />
4)We will have another go routine listening for incoming messages <br />
from the SendHeartBeatToNodes() go routine. Lets call this go-routine : **RebalanceContainersOnNodes()**<br />
Once a message is received from the go routine started in step #3, this go routine will inspect the<br />
node.db table to get the list of nodes which are offline.  With this information in hand, this go routine will then<br />
inspect the containers.db table to rebalance the affected containers.<br />
5)We will also have additional go routines to monitor the requests for addition/deletion to the cluster<br />
Lets call this go routine  **AddNodesToCluster()** and **DeleteNodesFromTheCluster()**<br />



Briefly the working of the REST API end points will be as follows:<br />
a)**/create  POST**<br />
Once there is a POST request to create a container with n number of replicas, the following workflow<br />
would kick in:<br />
1)The handler function(createContainerHandler) would get trigerred.<br />
2)The handler would invoke the container.CreateContainers() method in container.go<br />
In this method, we will use Mutexes from the sync.Mutex package or we will use a buffered channel with a capacity of 1<br />
to establish locks before the containers are spawned.<br />
This is to ensure that the spawning of containers is synchronized and container instances of only one image are being spawned<br />
at a time.<br />
3)Once the containers are spawned, the container.db table is updated.<br />

b)**/status GET**<br />
Once there is a GET request to get the status of the cluster, the Docker REST API can again be used to get the required information.<br />
However, if the required data can easily be fetched from the DB, we can do that as well<br />
The choice between fetching data from the DB vs using the Docker REST API is dependent upon the data we want to make available<br />
as part of the **/status** route <br />

**Spawning of containers**:<br />
1)All the containers can be spawned using the Docker remote API. The docker engine will have to be configured to accept REST connections.<br />
Using the Docker REST API allows us to spawn containers on remote systems as well. Hence, its the best choice to spawn containers on <br />
remote systems <br />
2)Instead of using different VM's acting as nodes, the other option is as follows:<br />
a)When the REST API server is started, spawn a fixed number of containers which have the **docker engine** installed in them<br />
These Docker containers would all be running on the same system and will act as nodes of the cluster<br />
b)Using the Docker REST API, run the **docker exec** command on each of these nodes as and when there is request for newer containers<br />
This would mean nesting of containers <br />

**Additional TODO:**<br />
-To interact with containers in container.go, make use of interfaces or rather program against interfaces.<br />
Using interfaces helps you extend your code for other container types easily and also makes the maintenance of the code easier.<br />
-Use interfaces to interact with sqllite db in db.go<br />
-Make the solution highly available by making the solution run behind a load balancer<br />
-Use a standard authentication package for authentication<br />
-Use a standard logging package<br />
-Write test cases by making use of GoLang's in built testing.T structure<br />
-Write benchmarking functions by making use of GoLang's in built benchmarking features<br />



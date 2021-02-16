package node

//A structure which represents the attributes of a cluster node.
type clusterNode struct {
	id           int
	nodeName     string
	nodeIpaddr   string
	nodeCapacity int64	//Represents the total capacity of the node
	nodeAvailCapacity int64	//Represents the available capacity of the node. 	 
	status 		 string //holds the status of whethere a node is online or offline
}

//ClusterNodes is a slice holding all the nodes in the cluster
//A slice has been used here since more items can be added to the slice
//as and when more nodes get added to the cluster.
type ClusterNodes []clusterNode

//Add is used to add more nodes to the cluster
//Return: error
func (c ClusterNodes) Add()  {

}

//This function will be used to periodically check the status of the cluster nodes
//and accordingly update the status of each node in the database as well as
//each node struct.
//Return:
func (c ClusterNodes) checkStatusOfClusterNodes() {

}
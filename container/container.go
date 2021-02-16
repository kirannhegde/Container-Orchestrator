package container

//A structure which represents the attributes of a container.
type container struct {
	containerRegistry string
	imageName         string
	imageVer          string
	requiredNumOfReplicas  int
	actualNumOfReplicas	int
}

type containers []container

func(c container) CreateContainer(cn container, numOfInstances int) {
	//Get the nodes in the cluster

	//Loop through the node.ClusterNodes slice
	//Check the capacity of each node
	//If capacity is available, create the container on that node.
	

	//Once the instances are created, update the details in the 
	//"containers" table

	//Update the node capacity in the "node" table against each node

}

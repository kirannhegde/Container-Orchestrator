package node

type clusterNode struct {
	id           int
	nodeName     string
	nodeIpaddr   string
	nodeCapacity int64
	status 		 string
}

type clusterNodes []clusterNode

package container

import (
	"encoding/json"
    "net/http"
	"bytes"
	"net/url"
    
)

//Image is a structure which represents the attributes of an image.
type Image struct {
	ImageNameTag string
}

//Container is a structure which represents the attributes of a container.
type Container struct {
	Image        //embedded type or inner type
	//imageTag          string
	RequiredNumOfReplicas  int
	actualNumOfReplicas	int
}



type containers []Container

//Function to pull the image from the specified repository
func (c Container) pullImage() (error) {
	imgNameTag := c.Image.ImageNameTag
	data := url.Values{}
	data.Set("fromImage", imgNameTag)
	response, err := http.PostForm("http://localhost:5555/images/create", data)
	defer response.Body.Close()
	return err
}

func(c Container) createContainer() (error) {
	postBody, err := json.Marshal(c.Image)
	if err != nil {
		return err
	}
	responseBody := bytes.NewBuffer(postBody)
	response, err := http.Post("http://localhost:5555/containers/create","application/json", responseBody)
	if err != nil {
		return err
	 }

	 defer response.Body.Close()

	 return nil

}

//CreateContainers is used to create containers on existing nodes
func(c Container) CreateContainers() (error) {
	//Get the nodes in the cluster

	//Loop through the node.ClusterNodes slice
	//Check the capacity of each node
	//If capacity is available, create the container on that node.
	

	//Once the instances are created, update the details in the 
	//"containers" table

	//Update the node capacity in the "node" table against each node

	err := c.pullImage()
	if err != nil {
		return err
	}

	for i:=0; i < c.RequiredNumOfReplicas; i++ {
		c.createContainer()
	}


	return nil
}

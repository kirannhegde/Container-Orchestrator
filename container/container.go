package container

//A structure which represents the attributes of a container.
type container struct {
	containerRegistry string
	imageName         string
	imageVer          string
}

type containers []container

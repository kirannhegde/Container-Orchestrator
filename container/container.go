package container

type container struct {
	containerRegistry string
	imageName         string
	imageVer          string
}

type containers []container

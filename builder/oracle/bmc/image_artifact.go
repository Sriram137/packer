package oracle

type ImageArtifact struct {
	displayName string
	imageId     string
	builderId   string
}

func (imageArtifact ImageArtifact) BuilderId() string {
	return imageArtifact.builderId
}

func (imageArtifact ImageArtifact) Id() string {
	return imageArtifact.imageId
}

func (imageArtifact ImageArtifact) String() string {
	return imageArtifact.displayName
}

func (imageArtifact ImageArtifact) Files() []string {
	return nil
}

func (imageArtifact ImageArtifact) State(string) interface{} {
	return nil
}

func (imageArtifact ImageArtifact) Destroy() error {
	return nil
}

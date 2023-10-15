package command

import "github.com/Foxcapades/Argonaut/v1/pkg/argo"

func Builder() argo.CommandBuilder {
	return &builder{}
}

type builder struct {
	description string
}

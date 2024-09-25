package types

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type ContainerRegistry interface {
	GetLatestImage(ctx context.Context, repositoryName string) (*types.ImageDetail, error)
	ImageExists(ctx context.Context, repositoryName string, imageTag string) error
}

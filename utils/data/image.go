package data

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func initializeECRClient(ctx context.Context, region string) (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("error loading AWS configuration: %w", err)
	}
	return ecr.NewFromConfig(cfg), nil
}

func GetLatestImage(ctx context.Context, repositoryName string, region string) (*types.ImageDetail, error) {
	svc, err := initializeECRClient(ctx, region)
	if err != nil {
		return nil, err
	}

	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repositoryName),
	}

	result, err := svc.DescribeImages(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error describing images: %w", err)
	}

	if len(result.ImageDetails) == 0 {
		return nil, fmt.Errorf("no images found in repository")
	}

	sort.Slice(result.ImageDetails, func(i, j int) bool {
		t1 := result.ImageDetails[i].ImagePushedAt
		t2 := result.ImageDetails[j].ImagePushedAt
		return t1.Before(*t2)
	})

	latestImage := result.ImageDetails[len(result.ImageDetails)-1]
	return &latestImage, nil
}

func ImageExists(ctx context.Context, repositoryName, imageTag, region string) error {
	svc, err := initializeECRClient(ctx, region)
	if err != nil {
		return err
	}

	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repositoryName),
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: aws.String(imageTag),
			},
		},
	}

	_, err = svc.DescribeImages(ctx, input)
	if err != nil {
		return fmt.Errorf("error describing images: %w", err)
	}

	return nil
}

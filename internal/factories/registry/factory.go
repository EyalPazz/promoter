package factories

import (
	"context"
	"fmt"
	"promoter/internal/types"
)

type RegistryFactory struct{}

func (f *RegistryFactory) InitializeRegistry(ctx context.Context, registryType, region string) (types.IContainerRegistry, error) {
	switch registryType {
	case "ecr":
		client, err := NewECRRegistryClient(ctx, region)
		if err != nil {
			return nil, err
		}
		return client, nil
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", registryType)
	}
}

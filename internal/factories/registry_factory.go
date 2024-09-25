package factories

import (
    "fmt"
    "context"
    "promoter/internal/data"
    "promoter/internal/types"
)

type RegistryFactory struct{}

func (f *RegistryFactory) InitializeRegistry(ctx context.Context, registryType string, region string )  (types.ContainerRegistry, error){
    switch registryType {
        case "ecr":
            client, err := data.NewECRRegistryClient(ctx, region)
            if err != nil {
            return nil, err
            }
            return client, nil
        default:
            return nil, fmt.Errorf("unsupported registry type: %s", registryType)
        }
}

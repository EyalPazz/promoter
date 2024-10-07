package commands

import (
	"context"
	"fmt"
	"promoter/internal/factories/registry"
	"promoter/internal/manipulations"
	"promoter/internal/types"
	"promoter/internal/utils"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd(cmd *cobra.Command, region string, services string, project string, env string, projectFile string) {
	passphrase, err := cmd.Flags().GetBool("passphrase")

	if region == "" {
		region = viper.GetString("region")
	}

	if region == "" {
		fmt.Println("Error: region must be specified either as flags or in the config file")
		return
	}

	if err = utils.RefreshRepo(passphrase); err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	serviceList, err := getServices(services, project, env, projectFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var changeLog []types.ServiceChanges

	// TODO: Think about the trade-offs in making this async
	for _, service := range serviceList {
		if err := processService(ctx, project, service, env, region, projectFile, &changeLog); err != nil {
			fmt.Println(err)
			if len(changeLog) > 0 {
				fmt.Println("Reverting Changes...")
				if err = manipulations.DiscardChanges(); err != nil {
					fmt.Println(err)
				}
			}
			return
		}

	}

	if len(changeLog) == 0 {
		fmt.Println("Nothing To Promote")
		return
	}

	commitMsg := utils.ComposeCommitMsg(&changeLog, env, project)

	if err := manipulations.HandleRepoActions(commitMsg, passphrase); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Success!")
}

func getServices(serviceStr string, project string, env string, projectFile string) ([]string, error) {
	var serviceList []string
	var err error

	if serviceStr == "" {
		serviceList, err = utils.GetServicesNames(project, env, projectFile)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving service names: %v", err)
		}
	} else {
		serviceList = strings.Split(serviceStr, ",")
	}

	return serviceList, nil
}

func processService(ctx context.Context, project string, service string, env string, region string, projectFile string, changeLog *[]types.ServiceChanges) error {
	repoName, err := utils.GetImageRepository(project, service, env, projectFile)
	if err != nil {
		return err
	}

	registryFactory := &factories.RegistryFactory{}

	// TODO: take type from image names after implementing more registries
	ecrClient, err := registryFactory.InitializeRegistry(ctx, "ecr", region)
	if err != nil {
		return err
	}

	latestImage, err := ecrClient.GetLatestImage(ctx, repoName)
	if err != nil {
		return err
	}
	tag := latestImage.ImageTags[len(latestImage.ImageTags)-1]

	didChange, err := manipulations.ChangeServiceTag(project, service, env, tag, projectFile)
	if err != nil {
		return err
	}

	if didChange {
		*changeLog = append(*changeLog, types.ServiceChanges{
			Name:   service,
			NewTag: tag,
		})
	}
	return nil
}

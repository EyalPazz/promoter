package commands

import (
	"context"
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/service"
	"promoter/internal/types"
	"promoter/internal/utils"
	"promoter/internal/utils/git"
	"strings"
)

type ProjectPromoter struct {
	env       string
	name      string
	Services  *[]service.Service
	ChangeLog *[]types.ServiceChanges
	Config    *utils.Config
	BaseCmd   types.IBaseCommand
}

func NewProjectPromoter(serviceString, env, name string) (*ProjectPromoter, error) {
	serviceList, err := getServices(serviceString, name, env)

	if err != nil {
		return nil, err
	}

	var services []service.Service

	config, err := utils.GetProjectConfig(name, env)
	if err != nil {
		return nil, err
	}

	for _, svc := range serviceList {
		service, err := service.NewService(name, svc, env, config)
		if err != nil {
			return nil, err
		}
		services = append(services, *service)
	}

	changeLog := &[]types.ServiceChanges{}

	project := ProjectPromoter{
		env,
		name,
		&services,
		changeLog,
		config,
		NewBaseCommand(env, name),
	}

	return &project, nil
}

func (p *ProjectPromoter) Process(tag, region string, interactive, passphrase bool) error {
	ctx := context.Background()

	for _, service := range *p.Services {
		if err := service.Process(ctx, tag, region, p.ChangeLog, interactive); err != nil {
			fmt.Println(err)
			if len(*p.ChangeLog) > 0 {
				fmt.Println(consts.RevertingChanges)
				workflow := git.DiscardGitFlow{BaseGitFlow: git.BaseGitFlow{}}
				if err = workflow.Execute(); err != nil {
					fmt.Println(err)
				}
			}
			return err
		}
	}

	if len(*p.ChangeLog) == 0 {
		return fmt.Errorf(consts.NothingToPromote)
	}

	return p.BaseCmd.Execute(passphrase, p.composeCommitTitle(), p.composeCommitBody())
}

func (p *ProjectPromoter) composeCommitTitle() string {
	return fmt.Sprintf("promotion(%s): %s \n", p.env, p.name)
}

func (p *ProjectPromoter) composeCommitBody() string {
	var msg string
	for _, change := range *p.ChangeLog {
		msg += fmt.Sprintf("changed %s to %s \n", change.Name, change.NewTag)
	}
	return msg
}

func getServices(serviceStr, project, env string) ([]string, error) {
	var serviceList []string
	var err error

	if serviceStr == consts.EmptyString {
		serviceList, err = utils.GetServicesNames(project, env)
		if err != nil {
			return nil, fmt.Errorf(consts.ErrorRetrievingServiceNames, err)
		}
	} else {
		serviceList = strings.Split(serviceStr, consts.Comma)
	}

	return serviceList, nil
}

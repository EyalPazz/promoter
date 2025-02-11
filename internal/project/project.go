package project

import (
	"context"
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/factories/gitprovider"
	"promoter/internal/service"
	"promoter/internal/types"
	"promoter/internal/utils"
	"promoter/internal/utils/git"
	"strings"
)

type Project struct {
	env       string
	name      string
	Services  *[]service.Service
	ChangeLog *[]types.ServiceChanges
	Config    *utils.Config
}

func NewProject(serviceString, env, name string) (*Project, error) {
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

	project := Project{
		env,
		name,
		&services,
		&[]types.ServiceChanges{},
		config,
	}

	return &project, nil
}

func (p *Project) Process(tag, region string, interactive, passphrase bool) error {
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

	return p.executeGitFlow(passphrase)
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

func (p *Project) executeGitFlow(passphrase bool) error {
	commitTitle := p.composeCommitTitle()
	commitBody := p.composeCommitBody()

	var workflow types.IGitFlow

	base_workflow := &git.BaseGitFlow{
		CommitMsg:  commitTitle + commitBody,
		Passphrase: passphrase,
	}

	if utils.ShouldCreatePR(p.env) {
		provider := gitprovider.GitProvider{}
		github := provider.GetProvider("github")
		workflow = &git.PRGitWorkflow{
			BaseGitFlow:  *base_workflow,
			GitProvider:  github,
			Title:        commitTitle,
			Body:         commitBody,
			ChangeBranch: utils.ComposeChangeBranch(p.name, p.env),
		}
	} else {
		workflow = base_workflow
	}

	if err := workflow.Execute(); err != nil {
		return err
	}

	return nil
}

func (p *Project) composeCommitTitle() string {
	return fmt.Sprintf("promotion(%s): %s \n", p.env, p.name)
}

func (p *Project) composeCommitBody() string {
	var msg string
	for _, change := range *p.ChangeLog {
		msg += fmt.Sprintf("changed %s to %s \n", change.Name, change.NewTag)
	}
	return msg
}

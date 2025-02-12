package commands

import (
	"promoter/internal/factories/gitprovider"
	"promoter/internal/types"
	"promoter/internal/utils"
	"promoter/internal/utils/git"
)

type BaseCommand struct {
	env  string
	name string
}

func NewBaseCommand(env, name string) *BaseCommand {
	return &BaseCommand{
		env, name,
	}
}

func (cmd *BaseCommand) Execute(passphrase bool, commitTitle, commitBody string) error {
	return cmd.executeGitFlow(passphrase, commitTitle, commitBody)
}

func (cmd *BaseCommand) executeGitFlow(passphrase bool, commitTitle, commitBody string) error {

	var workflow types.IGitFlow

	base_workflow := &git.BaseGitFlow{
		CommitMsg:  commitTitle + commitBody,
		Passphrase: passphrase,
	}

	if utils.ShouldCreatePR(cmd.env) {
		provider := gitprovider.GitProvider{}
		github := provider.GetProvider("github")
		workflow = &git.PRGitWorkflow{
			BaseGitFlow:  *base_workflow,
			GitProvider:  github,
			Title:        commitTitle,
			Body:         commitBody,
			ChangeBranch: utils.ComposeChangeBranch(cmd.name, cmd.env),
		}
	} else {
		workflow = base_workflow
	}

	if err := workflow.Execute(); err != nil {
		return err
	}

	return nil
}

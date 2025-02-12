package commands

import (
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/types"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type ProjectReverter struct {
	env     string
	name    string
	BaseCmd types.IBaseCommand
}

func NewProjectReverter(env, name string) *ProjectReverter {

	return &ProjectReverter{
		env,
		name,
		NewBaseCommand(env, name),
	}

}

func (p *ProjectReverter) Process(since int, passphrase bool) error {
	revs, err := utils.GetLatestRevisions(p.name, p.env, since)
	if err != nil {
		return err
	}

	commits := map[string]*object.Commit{}

	if len(revs) == 0 {
		return fmt.Errorf(consts.NoChangesRecorded)
	}

	var commitList []string
	for _, commit := range revs {
		commits[commit.Hash.String()] = commit
		commitList = append(commitList, commit.Hash.String())
	}

	var selected string

	prompt := &survey.Select{
		Message: consts.SelectRevision,
		Options: commitList,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return err
	}

	projectFile, err := utils.GetProjectFile(p.name, p.env, true)
	if err != nil {
		return err
	}

	file, err := utils.GetFileFromCommit(commits[selected], projectFile)
	if err != nil {
		return err
	}

	if err = utils.WriteToProjectFile(p.name, p.env, file); err != nil {
		return err
	}

	return p.BaseCmd.Execute(passphrase, p.composeCommitTitle(), p.composeCommitBody(selected))

}

func (p *ProjectReverter) composeCommitTitle() string {
	return fmt.Sprintf("demotion(%s): %s \n", p.env, p.name)
}

func (p *ProjectReverter) composeCommitBody(selected string) string {
	return fmt.Sprintf(consts.RevertingTo, p.name, p.env, selected)
}

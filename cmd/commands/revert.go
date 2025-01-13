package commands

import (
	"fmt"
	"promoter/internal/consts"
	"promoter/internal/manipulations"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func RevertProject(cmd *cobra.Command) {
	env, _ := cmd.Flags().GetString(consts.Env)
	passphrase, _ := cmd.Flags().GetBool(consts.Passphrase)
	project, _ := cmd.Flags().GetString(consts.Project)
	since, _ := cmd.Flags().GetInt(consts.Since)

	project, _, err := utils.ValidateProjectAttributes(project, consts.EmptyString)

	if project == consts.EmptyString {
		fmt.Print(err)
		return
	}

	revs, err := utils.GetLatestRevisions(project, env, since)
	if err != nil {
		fmt.Printf(consts.ErrorGettingPreviousRevisions, err)
		return
	}

	commits := map[string]*object.Commit{}

	if len(revs) == 0 {
		fmt.Println(consts.NoChangesRecorded)
		return
	}

	for _, commit := range revs {
		commits[commit.Hash.String()] = commit
	}

	var selected string

	prompt := &survey.Select{
		Message: consts.SelectRevision,
		Options: maps.Keys(commits),
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		fmt.Println(err.Error())
		return
	}

	projectFile, err := utils.GetProjectFile(project, env, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := utils.GetFileFromCommit(commits[selected], projectFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = utils.WriteToProjectFile(project, env, file); err != nil {
		fmt.Println(err)
		return
	}

	if err = manipulations.HandleRepoActions(fmt.Sprintf(consts.RevertingTo, project, env, selected), passphrase); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(consts.Success)
}

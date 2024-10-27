package commands

import (
	"fmt"
	"promoter/internal/manipulations"
	"promoter/internal/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func RevertProject(cmd *cobra.Command) {
	env, _ := cmd.Flags().GetString("env")
	passphrase, _ := cmd.Flags().GetBool("passphrase")
	project, _ := cmd.Flags().GetString("project")
	projectFile, _ := cmd.Root().PersistentFlags().GetString("project-file")
	since, _ := cmd.Flags().GetInt("since")

	if projectFile == "" && (env == "" || project == "") {
		fmt.Print("You Need to either provide both env and project flags, or the project-file path")
		return
	}

	revs, err := utils.GetLatestRevisions(project, env, since)
	if err != nil {
		fmt.Println("Error Getting Previous Revisions: ", err)
		return
	}

	commits := map[string]*object.Commit{}

	if len(revs) == 0 {
		fmt.Println("No changes recorded in given interval")
		return
	}

	for _, commit := range revs {
		commits[commit.Hash.String()] = commit
	}

	var selected string

	prompt := &survey.Select{
		Message: "Select a revision to revert to",
		Options: maps.Keys(commits),
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		fmt.Println(err.Error())
		return
	}

	if projectFile == "" {
		projectFile, err = utils.GetProjectFile(project, env, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file, err := utils.GetFileFromCommit(commits[selected], projectFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = utils.WriteToProjectFile(project, env, "", file); err != nil {
		fmt.Println(err)
		return
	}

	if err = manipulations.HandleRepoActions(fmt.Sprintf("Reverting %s(%s) to %s", project, env, selected), passphrase); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Success!")
}

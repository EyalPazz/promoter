package commands

import (
	"fmt"
	"promoter/internal/utils"

    "golang.org/x/exp/maps"
	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func RevertService(cmd *cobra.Command) {
	env, _ := cmd.Flags().GetString("env")
	project, _ := cmd.Flags().GetString("project")
	projectFile, _ := cmd.Root().PersistentFlags().GetString("project-file")

	if projectFile == "" && (env == "" || project == "") {
		fmt.Print("You Need to either provide both env and project flags, or the project-file path")
		return
	}

	revs, err := utils.GetLatestRevisions(project, env, 20)
	if err != nil {
		fmt.Println("Error Getting Previous Revisions: ", err)
		return
	}

    commits := map[string]*object.Commit{}

    for _ , commit := range revs {
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
    
    projectFilePath := projectFile
    if projectFilePath == "" {
        projectFilePath, err =  utils.GetProjectFile(project, env, true)
        if err != nil {
            fmt.Println(err)
        }
    }


	file, err := utils.GetFileFromCommit(commits[selected], projectFilePath)
	fmt.Println(file)
}

package gitprovider

import (
	"promoter/internal/types"
)

type GitProvider struct{}

func (provider *GitProvider) GetProvider(name string) types.IGitProvider {
	switch name {
	case "github":
		return NewGitHubClient()
	default:
		return NewGitHubClient()
	}
}

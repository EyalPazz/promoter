package types

type Config struct {
	GitName          string             `mapstructure:"git-name"`
	GitEmail         string             `mapstructure:"git-email"`
	SSHKey           string             `mapstructure:"ssh-key"`
	ManifestRepo     string             `mapstructure:"manifest-repo"`
	ManifestRepoRoot string             `mapstructure:"manifest-repo-root"`
	PullRequests     PullRequests       `mapstructure:"pullRequests"`
	Profiles         map[string]Profile `mapstructure:",remain"`
}

type Profile struct {
	ProjectName string `mapstructure:"project-name"`
	Region      string `mapstructure:"region"`
}

type PullRequests struct {
	Enabled    bool     `mapstructure:"enabled"`
	BaseBranch string   `mapstructure:"base-branch"`
	Org        string   `mapstructure:"org"`
	RepoName   string   `mapstructure:"repo-name"`
	Envs       []string `mapstructure:"envs"`
}

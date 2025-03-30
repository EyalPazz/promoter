# Promoter

Promoter facilitates easy promotion between environments using GitOps and Helm. It automates the update of image tags for services in different environments.

## Installation

```bash
curl -sL https://github.com/EyalPazz/promoter/raw/main/install.sh | sudo bash
```

Or, from the repository itself:

```bash
sudo make install
```

By default, this installs the binary into `/usr/local/bin`, which requires `sudo`. If you want to use a different prefix, you can set the `PREFIX` environment variable (the default is `/usr/local`, as `bin` is appended to the prefix):

```bash
curl -sL https://github.com/EyalPazz/promoter/raw/main/install.sh | PREFIX=$HOME/.local bash
```

Or, from the repository itself:

```bash
PREFIX=$HOME/.local make install
```

## Usage

### Promotion

To promote all services of a project to an environment with the latest image tag, use:

```bash
promoter  --project <project> --env <env>
```

To promote certain services of a project to an environment with the latest image tag, use:

```bash
promoter --services <services list, seperated by commas> --project <project> --env <env>
```

To promote a service to an environment with a specific image tag, use:

```bash
promoter --services <service_name> -tag <tag> --env <env> --project <project>
```

You can also use the -i (or --interactive) in order to promote interactively

### Reverting

To revert all services of a project to a certain revision

```bash
promoter revert --env production
```

** You can also use the project flag to override the .promoter.yaml config **

## Configuration

Create a configuration file named .promoter.yaml in your home directory (~/.promoter.yaml) with the following key-value pairs:git-name: <Your Git Username>

```yaml
git-email: <Your Git Email>
git-name: <Your Git Name>
manifest-repo: <Config Files Repo URL>
ssh-key: <Git SSH Key Path>
manifest-repo-root: < ** OPTIONAL ** For app of apps repos>
pullRequests:
  enabled: boolean (required)
  org: string
  base-branch: string
  repo-name: string
  envs: list<string> (all envs by default)

default:
  project-name: X
  region: Y
```

\*\* In order for the PR's to work, you should have an env variable named `GIT_PROVIDER_TOKEN` that stores your PAT, It Should Have Read and Write access to code and pull requests \*\*

You can also add another profile and use --profile with it's name

Your manifest project structure should follow this format:

```perl
argo/
  apps/
      ├── project1/
      │   └── env/
      │       └── values.yaml
      └── project2/
          └── env/
              └── values.yaml
```

## To-Do List

- [x] 1: Make Deployment Platform agnostic
- [x] 2: Write documentation
- [ ] 3: Add Tests
- [ ] 4: Expend to more providers and manifest repo structures

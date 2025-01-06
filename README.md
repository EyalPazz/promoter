# Promoter

Promoter facilitates easy promotion between environments using GitOps and Helm. It automates the update of image tags for services in different environments.

## Installation

```bash
curl -sL https://github.com/EyalPazz/promoter/raw/main/install.sh | bash
```

or

```bash
make install
```

## Usage

### Promotion

To promote all services of a project to an environment with the latest image tag, use:

```bash
promoter  --env <env>
```

To promote certain services of a project to an environment with the latest image tag, use:

```bash
promoter --services <services list, seperated by commas> --env <env>
```

To promote a service to an environment with a specific image tag, use:

```bash
promoter --services <service_name> -tag <tag> --env <env>
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

default:
  project-name: X
  region: Y
```

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
- [ ] 2: Add Tests
- [ ] 3: Expend to more providers and manifest repo structures
- [x] 4: Write documentation

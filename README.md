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

To promote a service to the production environment with the latest image tag, use:

```bash
promoter --project <project_name> --service <service_name> --env production
```

To promote a service to the production environment with a specific image tag, use:

```bash
promoter --project <project_name> --service <service_name> --env production
```

To promote all services of a project to the production environment with the latest image tag, use:

```bash
promoter --project <project_name> --env production
```

### Reverting

To revert all services of a project to a certain revision

```bash
promoter revert --project <project_name> --env production
```

## Configuration

Create a configuration file named .promoter.yaml in your home directory (~/.promoter.yaml) with the following key-value pairs:git-name: <Your Git Username>

```yaml
git-email: <Your Git Email>
git-name: <Your Git Name>
manifest-repo: <Config Files Repo URL>
ssh-key-path: <Git SSH Key Path>
region: <Region of Your Container Registry>
manifest-repo-root: < ** OPTIONAL ** For app of apps repos>
```

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
- [x] 4: Write documentation
- [ ] 2: Add Tests
- [ ] 3: Expend to more providers and manifest repo structures

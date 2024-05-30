# Promoter

Promoter facilitates easy promotion between environments using GitOps and Helm. It automates the update of image tags for services in different environments.

## Usage

To promote a service to the production environment with the latest image tag, use:

```bash
promoter --project <project_name> --service <service_name> --env production
```

To specify a custom image tag, use the --tag option:

```bash
promoter --project <project_name> --service <service_name> --env production --tag <image_tag>
```

## Configuration

Create a configuration file named .promoter.yaml in your home directory (~/.promoter.yaml) with the following key-value pairs:git-name: <Your Git Username>

```yaml
git-email: <Your Git Email>
manifest-repo-url: <K8s/Helm Files Repo URL>
ssh-key-path: <GitHub SSH Key Path>
region: <AWS Region of Your ECR Repo>
repository-name: <Name of Your ECR Repository>
```

Your manifest project structure should follow this format:

```perl
chart-repo/
  ├── project1/
  │   └── service/
  │       └── values-${env}.yaml
  └── project2/
      └── service/
          └── values-${env}.yaml
```

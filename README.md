# Slack bot for managing application versions on multi-env Kubernetes

This bot is designed to simplify the process of managing application versions on Kubernetes directly from your Slack workspace. Whether you need to check the current status of versions, view changes, promote to the next environment, or rollback to a previous version, this bot helps you.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [`list`](#list)
  - [`diff`](#diff)
  - [`promote`](#promote)
  - [`rollback`](#rollback)
  - [How it works](#how-it-works)
- [Requirements](#requirements)
  - [Functional](#functional)
  - [Non-Functional](#non-functional)
- [High-Level Design](#high-level-design)
- [Architecture Decision Records](#architecture-decision-records)

## Getting Started

### Prerequisites

Ensure that you have the following prerequisites before installing the Slack bot:

- Configured and running multi-environment Kubernetes clusters (for demo see [/demo/README.md](./demo/README.md))
- Installed and configured Flux to implement a GitOps approach to managing both infrastructure and applications. (for demo see [/demo/README.md](./demo/README.md))
- Setup and run slackbot service and add bot to slack channel ([instruction](/app/README.md))

### Installation

To use the Slackbot, follow these steps to install and set up the necessary credentials:

1. **Slack Bot Tokens:**
   - Create a Slack app on your workspace.
   - Generate a Bot Token (`SLACK_BOT_TOKEN`) and an App-level Token (`SLACK_APP_TOKEN`).
   
2. **GitHub Token:**
   - Obtain a GitHub Personal Access Token (`GITHUB_API_TOKEN`) with the necessary permissions.
   
3. **Service Configuration:**
   - Set the environment variables by adding the following lines to your shell configuration file (e.g., `.bashrc`, `.zshrc`):

     ```dotenv
     export SLACK_BOT_TOKEN=<xoxb-...>
     export SLACK_APP_TOKEN=<xapp-...>
     export GITHUB_API_TOKEN=<github_token>
     export SERVICE_CONFIG_FILE_PATH='path to config'
     ```
   
   - Configure the Slack bot using the provided [config](config.example.yaml) file. Indicate the file path in `SERVICE_CONFIG_FILE_PATH`. The configuration includes settings for clusters, Slack bot, and the GitHub repository.

## Usage bot

More details you can find in [`./app/README.md`](./app/README.md)

Once the Slackbot is set up, you can interact with it using various commands in the Slack channel:

### `list`

Use the `list` command to obtain an overview of the current version status in different environments:

```
@botName list
```

### `diff`

The `diff` command shows a list of changes needed to update the application version in a specific environment. Specify the target environment, e.g., `stage`:

```
@botName diff stage
```

### `promote`

Deploy your application to a specific cluster and environment using the `promote` command. Provide the service name, version, and target environment:

```
@botName promote podinfo@7.0.0 to prod
```

- `podinfo`: Service name (must match the name in the deployment manifest `metadata.name`).
- `7.0.0`: Version Helm package for the service.
- `prod`: Target environment/cluster (must match the name in the config bot file `SERVICE_CONFIG_FILE_PATH` clusters list).

### `rollback`

The `rollback` command reverts the application version update to the previous version. Provide the service name and target environment:

```
@botName rollback prometheus-deployment on stage
```

- `prometheus-deployment`: Service name (must match the name in the deployment manifest `metadata.name`).
- `stage`: Target environment/cluster (must match the name in the config bot file `SERVICE_CONFIG_FILE_PATH` clusters list).

### How It Works

The Slack bot accepts actions in the Slack channel, creating a job based on those actions. The job has access to all clusters and applications in different namespaces. The primary focus is deploying applications via `Deployment` Kubernetes resources, as the bot monitors deployments and tracks them automatically.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podinfo # application name inside service
```

- **podinfo** - we use it like application name inside service and track and detect applications git `deployment.GetName()`
- **@botName** - name for your bot, you setted it when created bot inside slack 

## Requirements

### Functional

1. The bot should have the capability to connect to the team channel in Slack.
2. Users should be able to invoke the `list` command to obtain the current status of application versions in different environments.
3. Users should have the ability to invoke the `diff` command to get a list of changes needed to update the application version in a specific environment.
4. Users should have the ability to invoke the `promote` command to perform an update of the application version to the next environment.
5. Users should have the ability to invoke the `rollback` command to revert the application version update to the previous version.

### Non-Functional

1. **Efficiency:** The system should be efficient and provide quick responses to user queries in Slack.
2. **Reliability:** The bot should be resilient to errors and ensure reliability in interactions with the GitOps system.
3. **Ease of Use:** Commands for version management should be simple and easy to understand for users.
4. **GitOps Integration:** The system should successfully integrate with GitOps infrastructure, providing automated deployment and update processes.
5. **Documentation:** The bot should have clear and comprehensive user documentation explaining the usage of each command and system capabilities.

## High-Level Design

Explore the high-level design and components of the Slack Version Management Bot in the [HLD Documentation](./HLD.md).

## Architecture Decision Records

For insights into the architectural decisions behind the Slack Version Management Bot, refer to the [ADR Documentation](./ADR.md).

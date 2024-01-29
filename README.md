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

Follow these steps to install and set up the bot:

## Usage bot

[Instruction](/app/README.md)

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

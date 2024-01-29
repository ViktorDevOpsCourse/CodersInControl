# Slackbot
Repo contains slack bot. You can add it in slack channel and work with your k8s clusters

----
## Dependencies
* Slack bot api and api, bot tokens
* GitHub flux infra repo
* k8s clusters with flux
* 
----
## Requirements
First you need create all credentials to work with external dependencies.

Setup below env variables 
```dotenv
export SLACK_BOT_TOKEN=<xoxb-...>
export SLACK_APP_TOKEN=<xapp-...>
export GITHUB_API_TOKEN=<github_token>
export SERVICE_CONFIG_FILE_PATH='path to config'
```

Configure slack bot [config](config.example.yaml) file which you can to indicate in SERVICE_CONFIG_FILE_PATH
```yaml
clusters: # list of clusters
  prod: # cluster name should be same like folder name in flux repo
    file: "/Users/viktor/.kube/config" # path to kube config for connect to cluster
  stage: # cluster name should be same like folder name in flux repo
    file: "/Users/viktor/.kube/config" # path to kube config for connect to cluster

bot: # related to slack to settings
  admins: # users who have permissions on user bot inside slack
    - "viktorzhabskiy"

repo: # config for github. Set work repo and branch
  owner: "ViktorDevOpsCourse"
  name: "flux-image-updates" # repo name
  branch: "main" # work branch, changes from this branch apply for clusters
```
----
## Run
To run it in local machine 
```go
go run main.go
```
---
## How to work with bot
### List
```dotenv
@botName list
```
# gobot

gobot is a chat bot build using Go

## Running gobot Locally

You must export the `SLACK_BOT_TOKEN` and `SLACK_API_TOKEN` environment variables before starting the app.

You can start gobot locally by running:

    `docker compose build && docker compose up`

### Configuration

`SLACK_BOT_TOKEN` a Hubot token

`SLACK_API_TOKEN` an API token with `usergroups:list` scope

You can create custom scripts and add them to the `src/bot` directory

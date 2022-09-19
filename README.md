# AWARE Helpers

This could be a nice CLI/TUI to do simple tasks like consistently generating telemetry data for testing.

## How To:


### Run

`go run cmd/aware/main.go`


### Build

Current System:

`go build -o dist/aware ./cmd/aware/main.go`


Other Architecture:

`env GOOS=windows GOARCH=amd64 go build -o dist/aware.exe ./cmd/aware/main.go`


### Format

Requires gofump: `go install mvdan.cc/gofumpt@latest`

`gofumpt -l -w .`


### Lint

Requires golangci-lint, install [here](https://golangci-lint.run/usage/install/)

`golangci-lint run -v`


## Insipration

- Bitwarden CLI (https://github.com/bitwarden/clients/tree/master/apps/cli)
- Jira CLI (https://github.com/ankitpokhrel/jira-cli)

## Core Libraries

- bubbletea
- viper
- cobra

## Implementation Thoughts

### Examples

- `aware init`
- `aware config set api http://localhost:3000`
- `aware config get api`
- `aware login`
    - Returns the JWT for the user to set as env variable
- `aware list activities`
    - Asks what activity type?
- `aware get activity {ID}`
- `aware generate device telemetry {ID}`
- `aware generate parameter telemetry {ID}`


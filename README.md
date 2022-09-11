# AWARE Helpers

This could be a nice CLI/TUI to do simple tasks like consistently generating telemetry data for testing
or stress testing, this will be written in Typescript, to take advantage of current knowledge and the ability
to reuse code.

## Insipration

- Bitwarden CLI (https://github.com/bitwarden/clients/tree/master/apps/cli)
- Jira CLI (https://github.com/ankitpokhrel/jira-cli)

## Potential Libraries

- commander.js
- inquirer
- chalk


## Implementation Thoughts

### Examples

- `aware config set api http://localhost:3000`
- `aware config get api`
- `aware login`
    - Returns the JWT for the user to set as env variable
- `aware list activities`
    - Asks what activity type?
- `aware get activity {ID}`
- `aware generate device telemetry {ID}`
- `aware generate parameter telemetry {ID}`


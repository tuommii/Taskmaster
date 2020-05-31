# Taskmaster

<p align="center">
	<img width=100% src="assets/taskmaster.png">
</p>

## Features

### Client
- [x] Terminal in raw ode
- [x] Autocompletion
- [x] History
- [x] Sends user input to server via TCP

| Key | Action |
|---------|---------|
| <kbd>Up</kbd> | Next history item |
| <kbd>Tab</kbd> | Next suggestion (autocomplete) |

#### Commands

| Command (alias) | Action |
|---------|---------|
| help (h) | Show help |
| status (st) | show job statuses |


### Server (daemon)

- Option to daemonize with `-d` flag
- Hot-reload config file by sending `SIGHUP`
- Clients are limited one at time
- If client exist, server remains listening for new client

## TODO
- Config validation
- config file via arg

## Notes

If u hot reload config, maybe scroll up to see was it reloaded

Send SIGHUP

`kill -1 4200`

Maybe use `/dev/null`

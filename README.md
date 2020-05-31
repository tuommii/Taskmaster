# Taskmaster

<p align="center">
	<img width=100% src="assets/taskmaster.png">
</p>

## Features

### Client
- [x] Terminal in raw mode
- [x] Autocompletion
- [x] History
- [x] Aliases
- [x] Sends user input to server via TCP

### Server (daemon)
- Option to daemonize with `-d` flag
- Hot-reload config file by sending `SIGHUP`
- Clients are limited one at time
- Server remains listening for a new client if client exits

## How To Use

### Keys
| Key | Action |
|---------|---------|
| <kbd>Up</kbd> | Next history item |
| <kbd>Tab</kbd> | Next suggestion (autocomplete) |

### Commands
| Command (alias) | Action |
|---------|---------|
| help (h) | Show help |
| status (st) | show job statuses |
| start (run) | start a job |
| stop (pause) | stop a job |
| restart | restart a job |
| reload | reload config file |

## TODO
- Config validation
- config file via arg
- check duplicates in config (reload also)

## Notes

If u hot reload config, maybe scroll up to see was it reloaded

Send SIGHUP

`kill -1 4200`

Maybe use `/dev/null`

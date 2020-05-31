<h1 align="center">
	Taskmaster
	<br>
	<img width=100% src="assets/taskmaster.png">
</h1>

## Features

### Client
- [x] UI with raw mode
- [x] Autocompletion
- [x] History
- [x] Aliases
- [x] Sends user input to server via TCP

### Server (daemon)
- [x] Option to daemonize with `-d` flag
- [x] Hot-reload config file by sending `SIGHUP`
- [x] Only one client at time is allowed
- [x] Server remains listening for a new client if client exits

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
- [ ] Logger
- [ ] Implement commands
- [ ] Get config file via arg
- [ ] Config validation
- [ ] Validate restarting processes
- [ ] Remove x-packages
- [ ] Mutex ?

## Notes

If u hot reload config, maybe scroll up to see was it reloaded

Send SIGHUP

`kill -1 4200`

Maybe use `/dev/null`

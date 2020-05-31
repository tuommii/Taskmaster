<h1 align="center">
	Taskmaster
	<br><br>
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
- [x] Hot-reload config file by sending `SIGHUP` signal
- [x] Only one client at time is allowed
- [x] Server remains listening for a new client if client exits

## How To Use

### Keys
| Key | Action |
|---------|---------|
| <kbd>**Up**</kbd> | Next history item |
| <kbd>**Tab**</kbd> | Next suggestion (autocomplete) |

### Commands
| Command (alias) | Action |
|---------|---------|
| <kbd>**help** (h)</kbd> | Show help |
| <kbd>**status** (st)</kbd> | show job statuses |
| <kbd>**start** (run)</kbd> | start a job |
| <kbd>**stop** (pause)</kbd> | stop a job |
| <kbd>**restart**</kbd> | restart a job |
| <kbd>**reload**</kbd> | reload config file |

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

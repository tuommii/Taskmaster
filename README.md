<h1 align="center">
	(WIP) Taskmaster
	<br><br>
	<img width=100% src="assets/taskmaster.png">
</h1>

## How it works

Screencast with [example config](#Example-config)

<p align="center">
	<img width=100% src="assets/taskmaster_1.gif">
</p>

## Features

### Client
- [x] UI with raw mode
- [x] History
- [x] *Autocompletion* for **commands**
- [x] *Autocompletion* for **task names**, server sends names to client when config is (re)loaded
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

### Example config
```json
{
	"realtime": {
		"command": "bash /home/miikka/projects/hive/taskmaster/assets/test.sh",
		"workingDir": "./assets",
		"startTime": 5,
		"numprocs": 1,
		"autostart": false,
		"autorestart": true,
		"stdout": "/tmp/realtime.log"
	},
	"failing": {
		"command": "notfound",
		"startTime": 2,
		"stopTime": 7,
		"numprocs": 1,
		"autostart": true,
		"autorestart": true,
		"stdout": "/tmp/failing.log"
	}
}
```

## TODO
- [x] Job names to client (suggestions)
- [ ] Refactor bad parts
- [ ] Logger
- [ ] Implement commands nicely
- [ ] Get config file via arg
- [ ] Config validation (check max lengths also)
- [ ] Validate restarting processes
- [ ] Remove x-packages
- [ ] Mutex ?
- [ ] Multiline support

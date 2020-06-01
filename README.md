<h1 align="center">
	Taskmaster
	<br><br>
	<img width=100% src="assets/taskmaster.png">
</h1>

## Features

### Client
- [x] UI with raw mode
- [x] History
- [x] *Autocompletion* for **commands**
- [x] *Autocompletion* for **task names** (server sends names to client when config is (re)loaded)
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

Config example, also in `assets` folder
```json
{
	"REALTIME": {
		"command": "bash test.sh",
		"workingDir": "/home/miikka/projects/hive/taskmaster/assets",
		"startTime": 3,
		"numprocs": 1,
		"autostart": true,
		"autorestart": true,
		"stdout": "/tmp/realtime"
	},
	"FAILING": {
		"command": "notexist",
		"workingDir": "/home/miikka/projects/hive/taskmaster/assets",
		"startTime": 2,
		"stopTime": 7,
		"retries": 3,
		"autostart": true,
		"autorestart": true,
		"stdout": "sadas/fdfd"
	},
	"exit value": {
		"command": "bash test3.sh",
		"workingDir": "/home/miikka/projects/hive/taskmaster/assets",
		"retries": 3,
		"autostart": true,
		"exitCodes": [0, 1, 2],
		"autorestart": true,
		"stdout": "sadas/fdfd"
	}
}
```

## TODO
- [ ] Job names to client (suggestions)
- [ ] Logger
- [ ] Implement commands via cli
- [ ] Get config file via arg
- [ ] Config validation (check max lengths also)
- [ ] Validate restarting processes
- [ ] Remove x-packages
- [ ] Mutex ?
- [ ] Multiline support

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
- UI with raw mode
- History
- *Autocompletion* for commands and **task names**
- Aliases
- Sends user input to server via TCP

### Server (daemon)
- Option to daemonize with `-d` flag
- Hot-reloads config file when receives `SIGHUP` signal
- Only one client at time is allowed
- Server remains listening for a new client if a client exits

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
| <kbd>**uptime**</kbd> | show uptime of process |
| <kbd>**fg**</kbd> | attach output to server stdout |
| <kbd>**bg**</kbd> | deattach output from server stdout |

### Example config
```json
{
	"realtime": {
		"command": "bash /home/miikka/projects/hive/taskmaster/assets/test.sh",
		"workingDir": "./assets",
		"startTime": 5,
		"instances": 1,
		"autostart": false,
		"autorestart": true,
		"stdout": "/tmp/realtime.log",
		"env": ["USER=Miikka"]
	},
	"failing": {
		"command": "notfound",
		"startTime": 2,
		"stopTime": 7,
		"instances": 1,
		"autostart": true,
		"autorestart": true,
		"stdout": "/tmp/failing.log"
	}
}
```

## TODO
- [x] Logger
- [ ] Flags to README
- [ ] Option for Ignore logging
- [ ] Compare config files when reloaded
- [x] Job names to client (suggestions)
- [x] Do process pool when loading jobs
- [x] Restart always/never/unexpected exit code
- [x] Refactor bad parts
- [x] Stop signal
- [x] Implement commands nicely [(result)](https://github.com/tuommii/Taskmaster/commit/41c430f1a634154d0f7831dc5001c05646a74bc4#diff-4c51d95b26718af4a2ac3bec33e54e49L128)
- [x] Get config file via arg
- [x] Config validation (check max lengths also)
- [x] Validate restarting processes
- [ ] Send fg to client
- [ ] Remove x-packages
- [ ] Mutex ?
- [ ] Windows support if i get windows machine
- [ ] Multiline support

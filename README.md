<h1 align="center">
	Taskmaster
	<br><br>
	<img width=100% src="assets/taskmaster.png">
</h1>

## How it works

Screencast (old) with example config

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
- Log to stdout with `-debug` flag
- Hot-reloads config file when receives `SIGHUP` signal
- Only one client at time is allowed
- Server remains listening for a new client if a client exits

## How To Use

Clone
```bash
git clone https://github.com/tuommii/Taskmaster.git
```

Build
```bash
cd Taskmaster/
go get ./...
make build
```

Run
```
./bin/server
./bin/client
```


### Keys
| Key | Action |
|---------|---------|
| <kbd>**Up**</kbd> | Next history item |
| <kbd>**Tab**</kbd> | Next suggestion (autocomplete) |

### Commands
| Command (alias) | Action |
|---------|---------|
| <kbd>**status** (st)</kbd> | show job statuses |
| <kbd>**reload**</kbd> | reload config file |
| <kbd>**start** (run) \<job_name\></kbd> | start a job |
| <kbd>**stop** (pause) \<job_name\></kbd> | stop a job |
| <kbd>**help** \<job_name\></kbd> | Show help |
| <kbd>**restart** \<job_name\></kbd> | restart a job |
| <kbd>**uptime** \<job_name\></kbd> | show uptime of process |
| <kbd>**fg** \<job_name\></kbd> | attach output to server stdout |
| <kbd>**bg** \<job_name\></kbd> | deattach output from server stdout |
| <kbd>**exit**</kbd> | quit |

### Example config
```json
{
	"realtime": {
		"command": "bash ./test.sh",
		"workingDir": "./assets",
		"startTime": 3,
		"instances": 1,
		"autostart": true,
		"autorestart": "unexpected",
		"stopSignal": "HUP",
		"env": ["USER=Miikka", "HOME=Hermanni"],
		"exitcodes": [0, 1]
	},
	"excode": {
		"command": "bash ./test3.sh",
		"workingDir": "./assets",
		"autostart": false,
		"startTime": 4,
		"autorestart": "false",
		"stopSignal": "HUP",
		"env": ["USER=Miikka", "HOME=Hermanni"],
		"exitcodes": [0, 1]
	}
}
```

## TODO
- [ ] Flags to README
- [ ] Validation rules to README
- [ ] SIGHUP -> Send new tasknames to client
- [ ] Mutex ?
- [ ] Remove x-packages (used for optional things) ?
- [ ] Windows support if i get windows machine
- [ ] Multiline support

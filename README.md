# Taskmaster

Todo app * Fun + UNIX = Taskmaster

## Requirements
- Only standard library allowed

## Features

- Client/Server, via TCP
- Selfmade CLI-interface
- fg command
- Hot-reload config file
- Only Stdlib is used
- clients are limited one at time
- if client exist, server remains listening for new client

## TODO
- Config validation

## Notes

If u hot reload config, maybe scroll up to see was it reloaded

Send SIGHUP

`kill -1 4200`

Maybe use `/dev/null`

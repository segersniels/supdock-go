# Supdock
[![CircleCI](https://circleci.com/gh/segersniels/supdock/tree/master.svg?style=shield)](https://circleci.com/gh/segersniels/supdock/tree/master)

What's Up, Dock(er)? A slightly more visual way to interact with the docker daemon. Supdock is a wrapper for the docker command meaning you can still use all of the other `docker` commands without issues.

<p align="center">
<img src="https://i.imgur.com/ATV0nP7.png" width="250">

## Why
Repetitive use of `docker ps`, `docker logs`, `docker stats` and `docker exec -ti` when troubleshooting  complex container setups can get chaotic. Supdock aims to optimize and speed up your workflow using docker.
Supdock also tries to expand on the basic docker commands adding custom commands for you to use.

<p align="center">
<img src="https://i.imgur.com/moY077k.gif" width="450">

## Installation
### Binary
Grab a binary from the [releases](https://github.com/segersniels/supdock-go/releases) page and move it into your desired `bin` (eg. `/usr/local/bin`) location.

### Go
```bash
go get -u github.com/segersniels/supdock
```

If you don't want to use `supdock` and `docker` separately you can just set an alias.

```bash
alias docker="supdock"
```

## Contributing
If you would like to see something added or you want to add something yourself feel free to create an issue or a pull request.

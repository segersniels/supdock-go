# Supdock
What's Up, Dock(er)? A slightly more visual way to interact with the docker daemon. Supdock is a wrapper for the docker command meaning you can still use all of the other `docker` commands without issues.

<p align="center">
<img src="https://i.imgur.com/ATV0nP7.png" width="250">

## Why
Repetitive use of `docker ps`, `docker logs`, `docker stats` and `docker exec -ti` when troubleshooting  complex container setups can get chaotic. Supdock aims to optimize and speed up your workflow using docker.

<p align="center">
<img src="https://i.imgur.com/moY077k.gif" width="450">

## Installation
Grab a binary from the [releases](https://github.com/segersniels/supdock-go/releases) page and move it into your desired `bin` location.

```bash
mv supdock_${VERSION}_${DISTRO}_amd64 /usr/local/bin/supdock
```

If you don't want to use `supdock` and `docker` separately you can just set an alias.

```bash
alias docker="supdock"
```

## Usage
```
Usage: supdock [options] [command]

Options:      
  -h, --help         output usage information

Commands:
  stop              Stop a running container
  start             Start a stopped container
  logs              See the logs of a container
  rm                Remove a container
  rmi               Remove an image
  prune             Remove stopped containers and dangling images
  stats             See the stats of a container
  ssh               SSH into a container
  history           See the history of an image
  env               See the environment variables of a running container
  latest, update    Update to the latest version of supdock
```

## Why a Go version?
As I have a running docker container count in my prompt, every new line basically executes `docker ps -q |wc -l |tr -d ' '`. But because I also have `docker` aliased to `supdock` the execution time for this command was getting noticeably slower as it was being executed through `node.js`. Which is not really known for it's optimal performance as it's single threaded.

```bash
supdock ps  0.33s user 0.07s system 103% cpu 0.390 total
supdock-go ps  0.06s user 0.02s system 91% cpu 0.091 total
```

## Contributing
If you would like to see something added or you want to add something yourself feel free to create an issue or a pull request.

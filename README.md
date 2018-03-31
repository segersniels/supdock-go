# Supdock
> Supdock in Golang is still under development. Installation section will be added when the package is deemed usable.

[![NPM](https://nodei.co/npm/supdock.png?mini=true)](https://npmjs.org/package/supdock)

What's Up, Dock(er)? A slightly more visual way to interact with the docker daemon. Supdock is a wrapper for the docker command meaning you can still use all of the other `docker` commands without issues.

<p align="center">
<img src="https://i.imgur.com/ATV0nP7.png" width="250">

## Why
Repetitive use of `docker ps`, `docker logs`, `docker stats` and `docker exec -ti` when troubleshooting  complex container setups can get chaotic. Supdock aims to optimize and speed up your workflow using docker.

<p align="center">
<img src="https://i.imgur.com/moY077k.gif" width="450">

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
```

## Contributing
If you would like to see something added or you want to add something yourself feel free to create an issue or a pull request.

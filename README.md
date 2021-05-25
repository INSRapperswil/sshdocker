# SSHDOCKER

## Requirements
- Docker

## Usage

Run a container of your choice.
```
$ docker run -d -i --name myubuntu ubuntu:latest bash
```

Start the sshdocker proxy to get SSH access to the container.
```
$ sshdocker --container-name myubuntu --shell /bin/bash --user customuser --password Super-8ecur3!
```

Connect to your Container.
```
$ ssh -p 2222 customuser@localhost
```

Additional settings `sshdocker --help`:
```
NAME:
   sshdocker - interactive ssh connection to container

USAGE:
   sshdocker [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --container-name value, -c value  Target container [$CONTAINER]
   --user value, -u value            User for authentication [$SSH_USER]
   --password value, -p value        Password for authentication [$SSH_PASSWORD]
   --host-key FILE, -k FILE          Host key FILE [$HOST_KEY_FILE]
   --shell value, -s value           Default shell (default: "/bin/sh") [$CONTAINER_SHELL]
   --port value                      Binding port (default: "2222") [$PORT]
   --help, -h                        show help (default: false)
```

## Docker Usage

```
docker run \
   -d -p 2222:2222 \
   -v /var/run/docker.sock:/var/run/docker.sock \
   -e CONTAINER=myubuntu \
   -e CONTAINER_SHELL=/bin/bash \
   -e SSH_USER=customuser \
   -e SSH_PASSWORD=Super-8ecur3! \
   insost/sshdocker:latest
```

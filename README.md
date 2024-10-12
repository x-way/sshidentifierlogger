# sshidentifierlogger

Log SSH identification strings. Collect information about which tools are used for SSH scanning/enumeration/bruteforcing against a machine.

## Installation

```
# go install github.com/x-way/sshidentifierlogger@latest
```

## Usage

```
# sshidentifierlogger -h
Usage of sshidentifierlogger:
  -d	debug mode, log to stdout
  -i string
    	Interface to read packets from (only supported on Linux) (default "eth0")
  -r string
    	Filename to read from, overrides -i
  -s string
    	Server IP, packets coming from this address are ignored (use commas to specify multiple IPs)

```

```
# sshidentifierlogger -d -i enp0s1
reading packets from interface enp0s1
{"src":"fdcc:433d:7bff:4567:de7b:4415:7778:5bc8","sport":41494,"dst":"fdcc:342e:7bff:9876::6","dport":22,"sshid":"SSH-2.0-OpenSSH_9.2p1 Debian-2+deb12u3","@timestamp":"2024-10-12T19:12:07Z","@version":1}
```

# narcotk-hosts

## What is narcotk-hosts?

narcotk-hosts is an simple hosts management application, the allows you to easily manage your cloud virtual machines and IOT devices.

## What do you mean manage?

- record VM or IOT device networking information
- add/delete/update information
- allow for devices to register them selves in narcotk-hosts
- provide a web api that allows devices to query their configuration
- provides a means to bootstrap VMs and IOT devices


## Features

- lightweight and simple to use
- tls/ssl encryption
- output in plain text or json
- can run stand alone or as a web service
- easy to run on osx, linux and windows.
- VMs and IOT devices can access the web api to get their boot strap script
- easy to run in a docker container
- can generate an old school hosts file
- IPv4 compatible
- IPv6 compatiblity to be added


## Example Uses

1. As a simple hosts file maintenance tool: you can run narcotk-hosts as a simple hosts file maintainer, adding and deleting hosts, and to generate a hosts file.
2. As a boot strapping tool: a VM or IOT device boots then runs ```\curl http://server.com/mac/de:ad:be:ef:ca:fe?script | bash ``` where narcotk-hosts provides a boot script that can configure your VM or IOT device.

## Installation

### Install from git

Requirements: go v1.9.1

```
git clone git@gitlab.com:narcotk/narcotk-hosts-2.git
cd narcotk-hosts-2
go get ./
go build main.go
```

### Install on Centos/Redhat/Fedora


### Install on Debian/Ubuntu


### Install on OSX


## Generating HTTPS Certificates and Keys

```bash
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

## Webapi URLS

| URL | Output |
|:--|:--|
| http://localhost:23000/hosts | lists all hosts |
| http://localhost:23000/hosts?json | list all hosts in json |
| http://localhost:23000/hosts?header | list all hosts with header |
| http://localhost:23000/host/**HOSTNAME** | print details for **HOSTNAME** |
| http://localhost:23000/host/**HOSTNAME**?json | print details for **HOSTNAME** in json |
| http://localhost:23000/host/**HOSTNAME**?header | print details for **HOSTNAME** with header |
| http://localhost:23000/host/**HOSTNAME**?script | print script for **HOSTNAME** |
| http://localhost:23000/networks | lists all networks |
| http://localhost:23000/networks?json | lists all networks in json |
| http://localhost:23000/network/**NETWORK_ID** | print details for **NETWORK_ID** |
| http://localhost:23000/network/**NETWORK_ID**?json | print details for **NETWORK_ID** in json |
| http://localhost:23000/ip/**IP** | print host details for **IP** |
| http://localhost:23000/ip/**IP**?json | print host details for **IP** in json |
| http://localhost:23000/mac/**MAC** | print host details for **MAC** |
| http://localhost:23000/mac/**MAC**?json | print host details for **MAC** in json |

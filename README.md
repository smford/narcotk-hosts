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
go build -o narcotk-hosts main.go
./narcotk-hosts --setupdb --database=./new-database-file.db
```

### Install on Centos/Redhat/Fedora


### Install on Debian/Ubuntu


### Install on OSX


### Use Docker



## Configuration

Configuration is possible three ways: using defaults, using a configuration file, or via command line arguments.  Configuration via environment variables will be added at a later date.

### First Run

When running narcotk-hosts for the first time, you need to run the below command to create a database file:

#### Create a default database file ./narcotk_hosts_all.db
```./narcotk-hosts --setupdb```

##### Create another database file /path/to/somefile.db
```./narcotk-hosts --setupdb --database=/path/to/somefile.db```


### Default Configuration

| Setting | Default | Details |
|:--|:--|:--|
| Database | ./narcotk_hosts_all.db | database file to use |
| EnableTLS | false | enable or disable TLS |
| HeaderFile | ./header.txt | display header file |
| JSON | false | print output as json |
| ListenPort | 23001 | port for narcotk-hosts to listen on |
| ListenIP | 127.0.0.1 | IP for narcotk-hosts to bind to |
| Scripts | ./scripts | directory of scripts |
| ShowHeader | false | show header, false by default |
| TLSCert | ./tls/server.crt | if EnableTLS true, use this TLS cert |
| TLSKey | ./tls/server.crt | if EnableTLS true, use this TLS key |
| Verbose | false | be verbose |


### Configuration File

The default configuration file (narco-hosts-config.json) is read from the same directory as the narcotk-hosts executable.

```
{
    "Database": "./narcotk_hosts_all.db",
    "EnableTLS": false,
    "HeaderFile": "./header.txt",
    "JSON": false,
    "ListenIP": "127.0.0.1",
    "ListenPort": "23001",
    "Scripts": "./scripts",
    "ShowHeader": false,
    "TLSCert": "./tls/server.crt",
    "TLSKey": "./tls/server.key",
    "Verbose": true
}
```


### Command Line Configuration Options

| Argument | Details | Example |
|:--|:--|:--|
| ```--configfile``` | Configuration File | --configfile=/path/to/file.yaml |
| ```--database``` | Database File | --database=/path/to/somefile.db |
| ```--json``` | Print output as JSON | --json |
| ```--listenip``` | IP for narcotk-hosts to bind to | --listenip=127.0.0.1 |
| ```--listenport``` | Port for narcotk-hosts to listen on | --listenport=23001 |
| ```--showheader``` | Show header | --showheader |


## Generating HTTPS Certificates and Keys

```bash
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```


## Web Service

narcotk-hosts can can as a command line tool or as a web service.

| Argument | Description |
|:--|:--|
| ```--startweb``` | starts the webservice using the configuration files setting for EnableTLS |
| ```--starthttp``` | force webservice to start using HTTP |
| ```--starthttps``` | force webservice to start using HTTPS |


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

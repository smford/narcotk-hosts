# narcotk-hosts

[![Build Status Travis](https://travis-ci.org/smford/narcotk-hosts.svg?branch=master)](https://travis-ci.org/smford/narcotk-hosts) [![Go Report Card](https://goreportcard.com/badge/github.com/smford/narcotk-hosts)](https://goreportcard.com/report/github.com/smford/narcotk-hosts) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


## What is narcotk-hosts?

narcotk-hosts is an simple hosts management application, the allows you to easily manage your virtual machines and IOT devices.  It is both a cli tool and has a web api, enabling for management via the command line and by allowing virtual machines or IOT devices access data stored within narcotk-hosts.


## What do you mean manage?

- record VM or IOT device networking information
- add/delete/update information
- allow for devices to self register in narcotk-hosts
- provide a web api that allows devices to query their configuration
- provides a means to bootstrap VMs and IOT devices


## Demo
A demo of narcotk-hosts running as the web api is available [https://hosts.narco.tk](https://hosts.narco.tk/)


## Features

- cli tool
- web api
- can run stand alone or as a web service
- lightweight and simple to use
- tls/ssl encryption
- basic auth
- oauth (soon)
- output in plain text or json
- easy to run on osx, linux and windows.
- self registration of VMs and IOT devices
- VMs and IOT devices can get host specific files, useful for them bootstrapping and configuring themselves
- easy to run in a docker container
- easy to run within heroku (free tier even) or other container services
- can generate an old school hosts file
- IPv4 and IPv6 compatible


## Example Uses

1. As a simple hosts file maintenance tool: you can run narcotk-hosts as a simple hosts file maintainer, adding and deleting hosts, and to generate a hosts file.
2. As a boot strapping tool: a VM or IOT device boots then runs `curl http://server.com:23000/mac/de:ad:be:ef:ca:fe?file=configure-system | bash` where narcotk-hosts provides a configure system script that can configure your VM or IOT device.


## Installation


### Install from git

Requirements:
- go v1.9.10
- dep v0.4.1

```
git clone git@github.com:smford/narcotk-hosts.git 
cd narcotk-hosts
dep ensure
go build -o narcotk-hosts main.go
./narcotk-hosts --setupdb --database=./new-database-file.db
```

### Install on Centos/Redhat/Fedora


### Install on Debian/Ubuntu


### Install on OSX



### Install on Heroku

To test, simply click the below button
[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

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
| DatabaseType | sqlite3 | database type to use (sqlite3 only supported at moment) |
| EnableTLS | false | enable or disable TLS |
| Files | ./files | directory of scripts |
| HeaderFile | ./header.txt | display header file |
| IndexFile | ./index.html | print index.html when user visits root web directory (http://server.com/) |
| JSON | false | print output as json |
| ListenPort | 23000 | port for narcotk-hosts to listen on |
| ListenIP | 127.0.0.1 | IP for narcotk-hosts to bind to |
| RegistrationKey | <blank> | Registration key to use when registering hosts, blank disables registration |
| ShowHeader | false | show header, false by default |
| TLSCert | ./tls/server.crt | if EnableTLS true, use this TLS cert |
| TLSKey | ./tls/server.crt | if EnableTLS true, use this TLS key |
| Verbose | false | be verbose |


### Configuration File

The default configuration file (narcotk-hosts-config.json) is read from the same directory as the narcotk-hosts executable.

```
{
    "Database": "./narcotk_hosts_all.db",
    "DatabaseType": "sqlite3",
    "EnableTLS": false,
    "Files": "./files",
    "HeaderFile": "./header.txt",
    "IndexFile": "./index.html",
    "JSON": false,
    "ListenIP": "127.0.0.1",
    "ListenPort": "23000",
    "RegistrationKey": "",
    "ShowHeader": false,
    "TLSCert": "./tls/server.crt",
    "TLSKey": "./tls/server.key",
    "Verbose": true
}
```


## Command Line Usage

### General Options
| Command | Description | Example |
|:--|:--|:--|
| `--displayconfig` | Prints out the applied configuration | |
| `--help` | Display help information |  |
| `--json` | Print output in json | |
| `--showheader` | Prepend headerfile to the output [default=false] | |
| `--version` | Display version | |


### Configuration and Database Options
| Command | Description | Example |
|:--|:--|:--|
| `--configfile` | Configuration file | --configfile=/path/to/file.yaml |
| `--database` | Database file | --database=/path/to/somefile.db |
| `--setupdb` | Setup a new blank database file | --setupdb  --database=./newfile.db |


### Host
| Command | Description | Example |
|:--|:--|:--|
| `--addhost` | Add a host (--addhost, --network and --ip are mandatory, the other params are optional) | --addhost=server-1-199.domain.com --network=192.168.1 --ip=192.168.1.13 --ipv6=::6 --short1=server-1-199 --short2=server --short3=serv --short4=ser --mac=de:ad:be:ef:ca:fe |
| `--delhost` | Delete a host (--delhost and --network are mandatory)| --delhost=server-1-200.domain.com --network=192.168.1 |
| `--host` | Display a host | --host=server1.domain.com |
| `--network` | Print all hosts in a network | --network=192.168.1 |
| `--showmac` | Show MAC addresses | --showmac |
| `--updatehost` | Update a host (--updatehost and --network are mandatory, other params are optional) | --updatehost=server-1-199.domain.com --network=192.168.1 --host=server-1-200.domain.com --newnetwork=192.168.1 --ip=192.168.1.200 --ipv6=::6 --short1=server-1-200 --short2=server --short3=serv --short4=ser --mac=de:ad:be:ef:ca:fe |


### Network
| Command | Description | Example |
|:--|:--|:--|
| `--addnetwork` | Add a new network | --addnetwork=192.168.2 --cidr=192.168.2.0/24 --desc="Management Network" |
| `--delnetwork` | Delete a network |--delnetwork=192.168.3 |
| `--listnetworks` | List all networks | --listnetworks |
| `--updatenetwork` | Update a network (--updatenetwork with one or more of --network, --cidr or --desc required) | --updatenetwork=192.168.2 --network=192.168.3 --cidr=192.168.3/24 --desc="3rd Management Network" |


### Web API
| Command | Description | Example |
|:--|:--|:--|
| `--listenip` | IP Address to listen on | --listenip=10.0.0.14 (ipv4 or ipv6)|
| `--listenport` | Port for webservice to listen on | --listenport=23000 |
| `--starthttp` | Start Web Service in foreground using HTTP | --starthttp |
| `--starthttps` | Start Web Service in foreground using HTTPS | --starthttps |
| `--startweb` | Start Web Service in foreground using config file EnableTLS setting | --startweb |


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
| `http://localhost:23000/host/HOSTNAME` | print details for **HOSTNAME** |
| `http://localhost:23000/host/HOSTNAME?file=motd` | download motd file for **HOSTNAME** |
| `http://localhost:23000/host/HOSTNAME?header=y` | print details for **HOSTNAME** with header |
| `http://localhost:23000/host/HOSTNAME?json=y` | print details for **HOSTNAME** in json |
| `http://localhost:23000/hosts` | lists all hosts |
| `http://localhost:23000/hosts?header=y` | list all hosts with header |
| `http://localhost:23000/hosts?json=y` | list all hosts in json |
| `http://localhost:23000/hosts?mac=y` | list all hosts with mac address |
| `http://localhost:23000/hosts?mac=y&header=y` | list all hosts with mac address and header|
| `http://localhost:23000/hosts/NETWORK_ID` | lists all hosts for a specific **NETWORK_ID** |
| `http://localhost:23000/hosts/NETWORK_ID?header=y` | list all hosts with header for a specific **NETWORK_ID** |
| `http://localhost:23000/hosts/NETWORK_ID?json=y` | list all hosts in json for a specific **NETWORK_ID** |
| `http://localhost:23000/hosts/NETWORK_ID?mac=y` | list all hosts with mac address for a specific **NETWORK_ID** |
| `http://localhost:23000/hosts/NETWORK_ID?mac=y&header=y` | list all hosts with header and mac address for a specific **NETWORK_ID**|
| `http://localhost:23000/ip/IP` | print host details for **IP** (either IPv4 or IPv6) |
| `http://localhost:23000/ip/IP?header=y` | print host details with header for **IP** (either IPv4 or IPv6) |
| `http://localhost:23000/ip/IP?json=y` | print host details for **IP** (either IPv4 or IPv6) in json |
| `http://localhost:23000/ip/IP?mac=y` | print host details with mac for **IP** (either IPv4 or IPv6) |
| `http://localhost:23000/mac/MAC` | print host details for **MAC** |
| `http://localhost:23000/mac/MAC?header=y` | print host details with header for **MAC** |
| `http://localhost:23000/mac/MAC?json=y` | print host details for **MAC** in json |
| `http://localhost:23000/networks` | lists all networks |
| `http://localhost:23000/networks?json=y` | lists all networks in json |
| `http://localhost:23000/network/NETWORK_ID` | print details for **NETWORK_ID** |
| `http://localhost:23000/network/NETWORK_ID?json=y` | print details for **NETWORK_ID** in json |


## Registration API

New hosts can be registered in to the database using the registration api call.  The registration api is only enabled when a RegistrationKey is set in the configuration file, to disable set RegistrationKey to "" (blank).

| Query | | Details | Example |
|:--|:--|:--|:--|
| key | **MANDATORY** | RegistrationKey (from configfile) | key=somepassword |
| fqdn | **MANDATORY** | hostname | fqdn=server1.domain.com |
| ip | **MANDATORY** | ip address | ip=10.10.1.67 |
| ipv6 | optional | ipv6 address | ipv6=::67 |
| nw | **MANDATORY** | network | nw=10.10.1 |
| s1 | optional | shortname 1 | s1=server1 |
| s2 | optional | shortname 2 | s2=s1 |
| s3 | optional | shortname 3 | s3=something1 |
| s4 | optional | shortname 4 | s4=somethingelse1 |
| mac | optional | mac address | mac=DE:AD:BE:EF:CA:FE |

### Examples
- ```curl https://server.com/register?key=password&fqdn=server1.domain.com&ip=10.10.1.67&nw=10.10.1```
- ```curl https://server.com/register?key=password&fqdn=server1.domain.com&ip=10.10.1.67&nw=10.10.1&mac=DE:AD:BE:EF:CA:FE&s1=server1&ipv6=::67```


## Files and Scripts
The web api can be used to present files and scripts back to a host.  These files and scripts can be used to do things such as configure the host.

In the configuration file, set the path to "Files" and place your files and scripts in to that directory.

Multiple files are possible, just pass the filename as a query.  Store the files in the files directory with the hostname has a prefix, for example:

| Filepath | Details | API Call Example |
| :-- | :-- | :-- |
| path/to/files/server1.something.com.example | example | `http://server.com:23000/host/server1.something.com?file=example` |
| path/to/files/server1.something.com.ifcfg-eth0 | ifcfg-eth0 | `http://server.com:23000/host/server1.something.com?file=ifcfg-eth0` |
| path/to/files/server1.something.com.motd | motd | `http://server.com:23000/host/server1.something.com?file=motd` |
| path/to/files/server2.something.com.ifcfg-eth1 | ifcfg-eth1 | `http://server.com:23000/host/server2.something.com?file=ifcfg-eth1` |
| path/to/files/server2.something.com.motd | motd | `http://server.com:23000/host/server2.something.com?file=motd` |

### Example Path Structure for Files and Scripts
![Example path structure for files and scripts](https://github.com/smford/narcotk-hosts/raw/master/images/files.png "Example path structure for files and scripts")


### Usage Examples
1. Download and run default file:

       `curl -sSL https://server.com:23000/host/server1.something.com?file=example | bash`


2. Download and install a machines MOTD:

       `wget http://server.com:23000/host/server1.something.com?file=motd -O /etc/motd`


## Bootstrapping a System

Assuming a vanilla machine boots and gets on the network via DHCP, this example will allow the system to configure itself by:

1. system finds its mac address
2. the system then queries the narcotk-hosts server for its hostname
3. the system then downloads and runs its script which then configures the rest of the system

```
MACADDRESS=$(ifconfig en1|grep ether|cut -f2 -d\ )
MYHOSTNAME=$(curl -s http://server.com:23000/mac/$MACADDRESS\?json | jq '.[].Hostname')
\curl -sSL https://server.com:23000/host/$MYHOSTNAME?file=example | bash
```


### CLI Examples
[![asciicast](https://asciinema.org/a/cgD7GgVVUKhnwgDuZJrroPSAo.png)](https://asciinema.org/a/cgD7GgVVUKhnwgDuZJrroPSAo)

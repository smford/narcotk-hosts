# Demo Site

### This is a demo of narcotk-hosts hosted on heroku, certain features are disabled.  To see further details including examples of how to use the cli please visit [Homepage](https://www.github.com/smford/narcotk-hosts)

---

# narcotk-hosts

[![Build Status Travis](https://travis-ci.org/smford/narcotk-hosts.svg?branch=master)](https://travis-ci.org/smford/narcotk-hosts) [![Go Report Card](https://goreportcard.com/badge/github.com/smford/narcotk-hosts)](https://goreportcard.com/report/github.com/smford/narcotk-hosts)


## What is narcotk-hosts?

narcotk-hosts is an simple hosts management application, the allows you to easily manage your virtual machines and IOT devices.  It is both a cli tool and has a web api, enabling for management via the command line and by allowing virtual machines or IOT devices access data stored within narcotk-hosts.


## What do you mean manage?

- record VM or IOT device networking information
- add/delete/update information
- allow for devices to self register in narcotk-hosts
- provide a web api that allows devices to query their configuration
- provides a means to bootstrap VMs and IOT devices


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
2. As a boot strapping tool: a VM or IOT device boots then runs ```\curl http://server.com:23000/mac/de:ad:be:ef:ca:fe?file=configure-system | bash ``` where narcotk-hosts provides a configure system script that can configure your VM or IOT device.

---

## Demo and API

### Hostname API
- Give details about a specific host
- Send files and scripts

| API | Example URL | Description |
|:--|:--|
| http://localhost:23000/host/HOSTNAME | [http://hosts.narco.tk/host/server1.something.com](http://hosts.narco.tk/host/server1.something.com) | print details for HOSTNAME |
| http://localhost:23000/host/HOSTNAME?file=motd | [http://hosts.narco.tk/host/server1.something.com?file=motd](http://hosts.narco.tk/host/server1.something.com?file=motd) | download motd file for HOSTNAME |
| http://localhost:23000/host/HOSTNAME?header=y | [http://hosts.narco.tk/host/server1.something.com?header=y](http://hosts.narco.tk/host/server1.something.com?header=y)| print details for HOSTNAME with header |
| http://localhost:23000/host/HOSTNAME?json=y | [http://hosts.narco.tk/host/server1.something.com?json=y](http://hosts.narco.tk/host/server1.something.com?json=y)| print details for HOSTNAME in json |


### Hosts API
- Present a list of hosts
- Generate a hosts file
- Print hosts within a specific network

| API | Example URL | Description |
|:--|:--|
| http://localhost:23000/hosts | [http://hosts.narco.tk/hosts](http://hosts.narco.tk/hosts) | lists all hosts |
| http://localhost:23000/hosts?header=y | [http://hosts.narco.tk/hosts?header=y](http://hosts.narco.tk/hosts?header=y) | list all hosts with header |
| http://localhost:23000/hosts?json=y | [http://hosts.narco.tk/hosts?json=y](http://hosts.narco.tk/hosts?json=y) | list all hosts in json |
| http://localhost:23000/hosts?mac=y | [http://hosts.narco.tk/hosts?mac=y](http://hosts.narco.tk/hosts?mac=y) | list all hosts with mac address |
| http://localhost:23000/hosts?mac=y&header=y | [http://hosts.narco.tk/hosts?mac=y&header=y](http://hosts.narco.tk/hosts?mac=y&header=y) | list all hosts with mac address and header|
| http://localhost:23000/hosts/NETWORK_ID | [http://hosts.narco.tk/hosts/192.168.2](http://hosts.narco.tk/hosts/192.168.2) | lists all hosts for a specific NETWORK_ID |
| http://localhost:23000/hosts/NETWORK_ID?header=y | [http://hosts.narco.tk/hosts/192.168.2?header=y](http://hosts.narco.tk/hosts/192.168.2?header=y) | list all hosts with header for a specific NETWORK_ID |
| http://localhost:23000/hosts/NETWORK_ID?json=y | [http://hosts.narco.tk/hosts/192.168.2?json=y](http://hosts.narco.tk/hosts/192.168.2?json=y) | list all hosts in json for a specific NETWORK_ID |
| http://localhost:23000/hosts/NETWORK_ID?mac=y | [http://hosts.narco.tk/hosts/192.168.2?mac=y](http://hosts.narco.tk/hosts/192.168.2?mac=y) | list all hosts with mac address for a specific NETWORK_ID |
| http://localhost:23000/hosts/NETWORK_ID?mac=y&header=y | [http://hosts.narco.tk/hosts/192.168.2?mac=ymac=y&header=y](http://hosts.narco.tk/hosts/192.168.2?mac=y&header=y) | list all hosts with header and mac address for a specific NETWORK_ID |


### IP API
- Give details about a specific host

| API | Example URL | Description |
|:--|:--|
| http://localhost:23000/ip/IP |[http://hosts.narco.tk/ip/192.168.1.1](http://hosts.narco.tk/ip/192.168.1.1) | print host details for IP (either IPv4 or IPv6) |
| http://localhost:23000/ip/IP?header=y | [http://hosts.narco.tk/ip/192.168.1.1?header=y](http://hosts.narco.tk/ip/192.168.1.1?header=y) | print host details with header for IP (either IPv4 or IPv6) |
| http://localhost:23000/ip/IP?json=y | [http://hosts.narco.tk/ip/192.168.1.1?json=y](http://hosts.narco.tk/ip/192.168.1.1/json=y) |print host details for IP (either IPv4 or IPv6) in json |
| http://localhost:23000/ip/IP?mac=y | [http://hosts.narco.tk/ip/192.168.1.1?mac=y](http://hosts.narco.tk/ip/192.168.1.1?mac=y) | print host details with mac for IP (either IPv4 or IPv6) |


### MAC API
- Give details about a specific host

| API | Example URL | Description |
|:--|:--|
| http://localhost:23000/mac/MAC | [http://hosts.narco.tk/mac/de:ad:be:ef:ca:fe](http://hosts.narco.tk/mac/de:ad:be:ef:ca:fe") | print host details for MAC |
| http://localhost:23000/mac/MAC?header=y | [http://hosts.narco.tk/mac/de:ad:be:ef:ca:fe?header=y](http://hosts.narco.tk/mac/de:ad:be:ef:ca:fe?header=y") | print host details with header for MAC |
| http://localhost:23000/mac/MAC?json=y | [http://hosts.narco.tk/mac/de:ad:be:ef:ca:fe?json=y](http://hosts.narco.tk/mac/de:ad:be:ef:ca:fe?json=y) | print host details for MAC in json |


### Networks API
- List all available networks

| API | Example URL | Description |
|:--|:--|
| http://localhost:23000/networks | [http://hosts.narco.tk/networks](http://hosts.narco.tk/networks) | lists all networks |
| http://localhost:23000/networks?json=y | [http://hosts.narco.tk/networks?json=y](http://hosts.narco.tk/networks?json=y)| lists all networks in json |


### Network Specific API
- Give details about a specific network

| API | Example URL | Description |
|:--|:--|
| http://localhost:23000/network/NETWORK_ID | [http://hosts.narco.tk/network/192.168.1](http://hosts.narco.tk/network/192.168.1) | print details for NETWORK_ID |
| http://localhost:23000/network/NETWORK_ID?json=y | [http://hosts.narco.tk/network/192.168.1?json=y](http://hosts.narco.tk/network/192.168.1?json=y) | print details for NETWORK_ID in json |

# narcotk-hosts


## Generating HTTPS Certificates and Keys

```bash
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

## URLS

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

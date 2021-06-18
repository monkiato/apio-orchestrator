<img src="docs/logo.svg" width="30%"/>

# Apio Orchestrator

[![Build Status](https://drone.monkiato.com/api/badges/monkiato/apio-orchestrator/status.svg)](https://drone.monkiato.com/monkiato/apio-orchestrator)
[![codecov](https://codecov.io/gh/monkiato/apio-orchestrator/branch/master/graph/badge.svg)](https://codecov.io/gh/monkiato/apio-orchestrator)
[![Go Report Card](https://goreportcard.com/badge/github.com/monkiato/apio-orchestrator)](https://goreportcard.com/report/github.com/monkiato/apio-orchestrator)

An orchestrator library for managing Apio nodes in docker.

Features:

 - Create new Apio nodes in docker
 - Auto-discovery feature through Traefik configuration
 - Using MongoDB as data storage
 - apio-orchestrator CLI for command line operation

## Sample files

Sample files and documentation available in [sample](sample) folder.

## CLI

To list available commands either run `apio-orchestrator` or `apio-orchestrator -h`
 
```
Apio orchestrator CLI is a management tool for Apio nodes running in docker with the ability
to handle new or existing Apio containers and modify their collections

Usage:
  apio-orchestrator [command]

Available Commands:
  addCollection    create a new collection definition for the specified node id
  addField         create a new field definition for the specified node id and collection
  create           create a new Apio node in docker
  help             Help about any command
  inspect          inspect node metadata
  ls               list all nodes
  remove           remove an existing docker container for the specified Apio node
  removeCollection remove an existing collection definition for the specified node id
  removeField      remove an existing field definition for the specified node id and collection
  start            start an existing docker container for the specified Apio node
  stop             stop an existing docker container for the specified Apio node
  update           update an existing docker container for the specified Apio node

Flags:
      --config string   config file absolute path
  -h, --help            help for apio-orchestrator
```

### Configuration

Default values:

	domain_name  = "localhost" //for Traefik configuration
	network_name = "apio-default"
	node_prefix  = "apio-"
	config_path  = "/var/lib/apio-orchestrator/"
	mongodb.host = "mongodb:27017"
	mongodb.name = "apio"


#### Change the config

To specify a custom config file, use the --config command line option. 
Or declare a `config.json` file in one of the following paths:

 - /etc/apio-orchestrator/
 - $HOME/.apio-orchestrator/
 - ./ (current working directory)
  
e.g. for config.json

    {
        "domain_name": "mydomain.com",
        "mongodb": {
            "host": "my-mongo-container:27017"
        }
    }

#### config.json properties

 - **domain_name**: the domain name used for auto-discovery feature through Traefik.
 e.g. for a domain name `mydomain.com` and a node called `library` the url to access collections
 available in this node is `library.mydomain.com`
 - **network_name**: an existing docker network name that will be use for every apio node. 
 The MongoDB container must be available in the same network.
 - **node_prefix**: prefix used for docker container names created by the orchestrator in order to
 prevent collision with other containers using same names.
 - **config_path**: root path used for internal config stored by the orchestrator. **Note**: if the path is changed
 the orchestrator will lost reference to existing nodes and their info, ensure this folder is moved manually to the new path.
 - **mongodb.host**: container name and port for the existing MongoDB instance
 - **mongodb.name**: name of the DB used to store node collections data
 
## Build Project
 
To build the project and generate the api-orchestrator CLI run:
 
     go build -o bin/apio-orchestrator .
 
 # 'Apio Orchestrator' Logo

<div>Icons made by <a href="https://www.flaticon.com/authors/freepik" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
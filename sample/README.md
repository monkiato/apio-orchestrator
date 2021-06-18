# Apio Orchestrator Example

### Checklist

1) Ensure Docker is up & running
2) Modify default values through `config.json` file (sample file available in this folder modifies domain name and base config path)
3) Create Docker network for Traefik and apio nodes (`docker network create apio-default`). If network name is different than default, it must be declared in `config.json`
4) Ensure Traefik and MongoDB are running under the same network (check [docker-compose.yml](docker-compose.yml) sample file)
   

### Run test

The environment is ready to run the orchestrator:

 - Create API node called `notebook` using [notebook.json](notebook.json) manifest. A collection called `people` will be available in this API.

    `api-orchestrator create notebook -m notebook.json`


 - Remove node

    `apio-orchestrator remove notebook`


 - Add collection and fields

    ```
    apio-orchestrator createCollection notebook countries
    apio-orchestrator createField notebook countries name
    apio-orchestrator createField notebook countries country_code
    ```

 - Inspect API

    `apio-orchestrator inspect notebook`
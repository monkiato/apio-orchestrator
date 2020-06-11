package config

const (
	//DefaultDomainName default domain name used if no other domain name was declared in CLI configuration
	DefaultDomainName  = "localhost"
	//DefaultNetworkName default network name used if no other network name was declared in CLI configuration
	DefaultNetworkName = "apio-default"
	//DefaultNodePrefix default node container prefix used if no other prefix was declared in CLI configuration
	DefaultNodePrefix  = "apio-"
	//DefaultConfigPath default path used to store node metadata and manifest required for the orchestration.
	//this path can be replaced through the CLI configuration
	DefaultConfigPath  = "/var/lib/apio-orchestrator/"
	//DefaultMongoDbHost default hostname used for mongodb connection used for node collections data
	DefaultMongoDbHost = "mongodb:27017"
	//DefaultMongoDbName default mongo database name used for mongodb connection used for node collections data
	DefaultMongoDbName = "apio"
	//NodeFolder subfolder used to store nodes metadata and manifests
	NodeFolder         = "nodes/"
	//ApiDockerImage is the latest apio docker image available used for orchestrator
	ApiDockerImage     = "docker.pkg.github.com/monkiato/apio/apio:0.5-alpha"
)

package proxy

//GetTraefikLabels create essential Traefik labels required for exposing a node into a specific public domain path
func GetTraefikLabels(
	nodeName string,
	hostname string,
	networkName string,
	exposePort int) map[string]string {
	return map[string]string{
		"traefik.http.routers." + nodeName + ".rule": "Host(`" + hostname + "`)",
		"traefik.docker.network":                     networkName,
		"traefik.enable":                             "true",
		"traefik.http.services." + nodeName + ".loadbalancer.server.port": string(exposePort),
	}
}

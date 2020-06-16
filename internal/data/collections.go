package data

// CollectionDefinition contains the main name structure for a rest API collection
type CollectionDefinition struct {
	Name   string            `json:"name"`
	Fields map[string]string `json:"fields"`
}

package models

type APIData struct {
	Questions  []Question  `json:"questions"`
	Countries  []Country   `json:"countries"`
	Regions    []Region    `json:"regions"`
	Continents []Continent `json:"continents"`
	Manifests  Manifests   `json:"manifests"`
}

type Manifests struct {
	Geography        ManifestInfo `json:"geography"`
	GeneralKnowledge ManifestInfo `json:"general_knowledge"`
}

type ManifestInfo struct {
	Version   string `json:"version"`
	UpdatedAt string `json:"updated_at"`
}

type APIRootResponse struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Datasets    map[string]interface{} `json:"datasets"`
	Endpoints   []EndpointInfo         `json:"endpoints"`
	Stats       map[string]int         `json:"stats"`
}

type EndpointInfo struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

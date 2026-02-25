package types

// Provider is a swap provider and its supported chains (GET /providers).
type Provider struct {
	Name               string   `json:"name"`
	Provider           string   `json:"provider"`
	DisplayName        string   `json:"displayName,omitempty"`
	DisplayNameLong    string   `json:"displayNameLong,omitempty"`
	Keywords           []string `json:"keywords,omitempty"`
	Count              int      `json:"count"`
	LogoURI            string   `json:"logoURI,omitempty"`
	URL                string   `json:"url,omitempty"`
	EnabledChainIds    []string `json:"enabledChainIds,omitempty"`
	SupportedActions   []string `json:"supportedActions,omitempty"`
	SupportedChainIds  []string `json:"supportedChainIds,omitempty"`
}

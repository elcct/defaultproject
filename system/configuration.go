package system

type ConfigurationDatabase struct {
	Hosts string `json:"host"`
	Database string `json:"database"`
}

type Configuration struct {
	Secret   string `json:"secret"`
	PublicPath string `json:"public_path"`
	TemplatePath string `json:"template_path"`
	Database ConfigurationDatabase
}

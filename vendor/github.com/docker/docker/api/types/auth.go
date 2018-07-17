package types 


type AuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`

	
	
	
	Email string `json:"email,omitempty"`

	ServerAddress string `json:"serveraddress,omitempty"`

	
	
	IdentityToken string `json:"identitytoken,omitempty"`

	
	RegistryToken string `json:"registrytoken,omitempty"`
}

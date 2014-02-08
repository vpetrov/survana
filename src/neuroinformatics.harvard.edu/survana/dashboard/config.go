package dashboard

type AuthConfig struct {
    Type string `json:"type,omitempty"`
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
}

type Config struct {
    Authentication *AuthConfig   `json:"authentication,omitempty"`
    AllowRegistration bool `json:"allow_registration,omitempty"`
}

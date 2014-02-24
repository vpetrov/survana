package auth

type Config struct {
    Type string `json:"type,omitempty"`
    AllowRegistration bool `json:"allow_registration,omitempty"`
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
}

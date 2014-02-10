package auth

type Config struct {
    Type string `json:"type,omitempty"`
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
}

package dashboard

import (
        "github.com/vpetrov/perfect/auth"
       )

type Config struct {
    Authentication *auth.Config `json:"authentication,omitempty"`
}

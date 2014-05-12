package dashboard

import (
        "github.com/vpetrov/perfect/auth"
       )

type Config struct {
    Authentication *auth.Config `json:"authentication,omitempty"`
    StoreUrl       string       `json:"store_url,omitempty"`
}

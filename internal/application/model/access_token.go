package model

type (
	AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	AccessToken struct {
		Value       string
		ExpiresIn   int64
		RequestedAt int64
	}
)

package model

type (
	AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	AccessToken struct {
		User        string `json:"user"`
		Value       string `json:"value"`
		ExpiresIn   int64  `json:"expires_in"`
		RequestedAt int64  `json:"requested_at"`
	}
)

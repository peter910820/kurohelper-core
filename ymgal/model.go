package ymgal

type (
	config struct {
		Endpoint     string
		ClientID     string
		ClientSecret string
	}
)

type (
	basicResp[T any] struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    T    `json:"data"`
	}

	tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}
)

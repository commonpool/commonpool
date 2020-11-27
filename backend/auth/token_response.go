package auth

// TokenResponse response from the oidc provider
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int64  `json:"notBeforePolicy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

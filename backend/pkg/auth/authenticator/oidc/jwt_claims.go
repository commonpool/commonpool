package oidc

// JwtClaims Claims part of the oidc response
type JwtClaims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
	// Audience          []string `json:"aud"`
	Expiration        int64  `json:"exp"`
	IssuedAt          int64  `json:"iat"`
	Id                string `json:"jti"`
	Type              string `json:"typ"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	EmailVerified     bool   `json:"email_verified"`
}

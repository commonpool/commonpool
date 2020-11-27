package auth

// RedirectResponseMeta API response when redirecting
type RedirectResponseMeta struct {
	RedirectTo string `json:"redirectTo"`
}

// RedirectResponse API response container
type RedirectResponse struct {
	Meta RedirectResponseMeta `json:"meta"`
}

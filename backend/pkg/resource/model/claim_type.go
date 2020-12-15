package model

type ClaimType string

const (
	OwnershipClaim ClaimType = "owner"
	ManagerClaim   ClaimType = "manager"
	ViewerClaim    ClaimType = "viewer"
)

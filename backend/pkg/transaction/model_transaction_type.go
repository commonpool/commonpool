package transaction

type Type string

const (
	ResourceSharedWithGroup  Type = "resource_shared_with_group"
	ResourceRemovedFromGroup Type = "resource_removed_from_group"
	ServiceProvided          Type = "service_provided"
	ResourceBorrowed         Type = "resource_borrowed"
	ResourceReturned         Type = "resource_returned"
	ResourceTaken            Type = "resource_taken"
	TimeCreditsExchanged     Type = "time_credits_exchanged"
)

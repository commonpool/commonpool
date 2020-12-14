package transaction

type TransactionType string

const (
	ResourceSharedWithGroup  TransactionType = "resource_shared_with_group"
	ResourceRemovedFromGroup TransactionType = "resource_removed_from_group"
	ServiceProvided          TransactionType = "service_provided"
	ResourceBorrowed         TransactionType = "resource_borrowed"
	ResourceReturned         TransactionType = "resource_returned"
	ResourceTaken            TransactionType = "resource_taken"
	TimeCreditsExchanged     TransactionType = "time_credits_exchanged"
)

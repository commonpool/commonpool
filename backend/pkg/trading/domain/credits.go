package domain

import "time"

type Transaction struct {
	credit      time.Duration
	debit       time.Duration
	description string
	accountID   string
	accountType string
}

type Entry struct {
	Transactions []Transaction
}

type Journal struct {
	Entries []Entry
}

package models

type User struct {
	ID       string
	Email    string
	Password string
}

type UserBalance struct {
	CurrencyCode string
	Actual       float64
	Reserved     float64
}

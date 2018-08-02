package model

type User struct {
	Id		uint64	`json:"user_id"`
	Amount	uint64	`json:"amount"`
}

type Transfer struct {
	FromUser	uint64	`json:"from_user"`
	ToUser		uint64	`json:"to_user"`
	Amount		uint64	`json:"amount"`
}

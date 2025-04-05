package domain

type Payment struct {
	Id          string  `json:"id"`
	PaymentDate int64   `json:"paymentDate"`
	Sender      string  `json:"sender,omitempty"`
	Recipient   string  `json:"recipient,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
}

package events

//ApplicationCreated event
type ApplicationCreated struct {
	LoanAmount int `json:"loan_amount"`
}

//LoanAmountChanged event
type LoanAmountChanged struct {
	LoanAmount int `json:"loan_amount"`
}

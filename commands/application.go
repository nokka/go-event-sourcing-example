package commands

import "github.com/mishudark/eventhus"

//CreateApplication assigned to an applicant.
type CreateApplication struct {
	eventhus.BaseCommand
	LoanAmount int
}

//ChangeLoanAmount of an application.
type ChangeLoanAmount struct {
	eventhus.BaseCommand
	LoanAmount int
}

package domain

import (
	"errors"

	"github.com/mishudark/eventhus"
	"github.com/nokka/go-event-sourcing-example/commands"
	"github.com/nokka/go-event-sourcing-example/events"
)

//ErrInvalidLoanAmount when the loan amount is invalid.
var ErrInvalidLoanAmount = errors.New("The given loan amount is not within the accepted range.")

// Application is the heart of the domain model.
type Application struct {
	eventhus.BaseAggregate
	LoanAmount int
}

// ApplyChange is called when ever an event is triggered on an application,
// it will handle the incoming event and change the application accordingly.
func (a *Application) ApplyChange(event eventhus.Event) {
	switch e := event.Data.(type) {
	case *events.ApplicationCreated:
		a.ID = event.AggregateID
		a.LoanAmount = e.LoanAmount
	case *events.LoanAmountChanged:
		a.LoanAmount = e.LoanAmount
	}
}

// HandleCommand takes care of an incoming command and creates an event
// based on that command.
func (a *Application) HandleCommand(command eventhus.Command) error {
	// Create a new event.
	event := eventhus.Event{
		AggregateID:   a.ID,
		AggregateType: "Application",
	}

	switch c := command.(type) {
	// An application is being created, we'll set the data from
	// the command onto the event.
	case commands.CreateApplication:
		event.AggregateID = c.AggregateID
		event.Data = &events.ApplicationCreated{
			LoanAmount: c.LoanAmount,
		}
	// Loan amount should be changed, we'll set the loan amount
	// from the command onto the event.
	case commands.ChangeLoanAmount:
		if c.LoanAmount < 0 {
			return ErrInvalidLoanAmount
		}

		event.Data = &events.LoanAmountChanged{
			LoanAmount: c.LoanAmount,
		}
	}

	// Commit the changes to the aggregate.
	a.BaseAggregate.ApplyChangeHelper(a, event, true)

	return nil
}

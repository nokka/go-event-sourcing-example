package main

import (
	"log"
	"time"

	"github.com/mishudark/eventhus"
	async "github.com/mishudark/eventhus/commandbus"
	"github.com/mishudark/eventhus/commandhandler/basic"
	"github.com/mishudark/eventhus/eventbus/rabbitmq"
	"github.com/mishudark/eventhus/eventstore/mongo"
	"github.com/mishudark/eventhus/utils"
	"github.com/nokka/go-event-sourcing-example/commands"
	"github.com/nokka/go-event-sourcing-example/domain"
	"github.com/nokka/go-event-sourcing-example/events"
)

func main() {

	end := make(chan bool)

	// Register events.
	reg := eventhus.NewEventRegister()
	reg.Set(events.ApplicationCreated{})
	reg.Set(events.LoanAmountChanged{})

	// Setup event store.
	eventstore, err := mongo.NewClient("localhost", 27017, "application-store")
	if err != nil {
		log.Fatal(err)
	}

	// Setup event bus, we'll use this to emit and register events on commands.
	rabbit, err := rabbitmq.NewClient("guest", "guest", "localhost", 5672)
	if err != nil {
		log.Fatal(err)
	}

	// Repository.
	repository := eventhus.NewRepository(eventstore, rabbit)

	// Setup command handler.
	commandRegister := eventhus.NewCommandRegister()
	commandHandler := basic.NewCommandHandler(repository, &domain.Application{}, "test-domain", "application")

	// Register commands to the command handler.
	commandRegister.Add(commands.CreateApplication{}, commandHandler)
	commandRegister.Add(commands.ChangeLoanAmount{}, commandHandler)

	// Create a new command bus with the handler.
	bus := async.NewBus(commandRegister, 30)

	// Let's create some applications.
	for i := 0; i < 3; i++ {
		go func() {
			uuid, err := utils.UUID()
			if err != nil {
				return
			}

			// Create an application.
			var application commands.CreateApplication
			application.AggregateID = uuid
			application.LoanAmount = 50000

			bus.HandleCommand(application)

			// Perform a change on the loan amount.
			time.Sleep(time.Millisecond * 100)
			change := commands.ChangeLoanAmount{
				LoanAmount: 300000,
			}

			change.AggregateID = uuid
			change.Version = 1

			bus.HandleCommand(change)
		}()
	}

	<-end
}

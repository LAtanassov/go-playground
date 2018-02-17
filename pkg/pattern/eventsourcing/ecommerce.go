package eventsourcing

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/altairsix/eventsource"
	"github.com/altairsix/eventsource/dynamodbstore"
)

// Events --------------

// OrderCreated Event
// a fact expressed in the past tense
type OrderCreated struct {
	eventsource.Model
}

// OrderApproved Event
type OrderApproved struct {
	eventsource.Model
}

// OrderShipped Event
type OrderShipped struct {
	eventsource.Model
}

// Commands --------------

// CreaterOrder Command
type CreateOrder struct {
	eventsource.CommandModel
}

// ApproveOrder Command
type ApproveOrder struct {
	eventsource.CommandModel
}

// Aggregates --------------

// Order is an Aggregate which apply Events
type Order struct {
	ID        string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
}

const (
	StateCreated  = "created"
	StateApproved = "approved"
	StateShipped  = "shipped"
)

// On an incoming Event apply updates to Aggregate
// current state
func (o *Order) On(event eventsource.Event) error {
	switch v := event.(type) {
	case *OrderCreated:
		o.CreatedAt = v.At
		o.State = StateCreated
	case *OrderApproved:
		o.State = StateApproved
	case *OrderShipped:
		o.State = StateShipped
	default:
		return fmt.Errorf("unhandeld event, %v", v)
	}

	o.ID = event.AggregateID()
	o.Version = event.EventVersion()
	o.UpdatedAt = event.EventAt()

	return nil
}

// Apply generates Events from commands
func (o *Order) Apply(ctx context.Context, command eventsource.Command) ([]eventsource.Event, error) {
	switch v := command.(type) {
	case *CreateOrder:
		orderCreated := &OrderCreated{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.Version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderCreated}, nil
	case *ApproveOrder:
		if o.State != StateCreated {
			return nil, fmt.Errorf("only %v orders may be approved", StateCreated)
		}
		orderApproved := &OrderApproved{
			Model: eventsource.Model{ID: v.AggregateID(), Version: o.Version + 1, At: time.Now()},
		}
		return []eventsource.Event{orderApproved}, nil
	default:
		return nil, fmt.Errorf("unhandled command %v", v)
	}
}

// How to generate the Events
// Command -> Command Handler -> Events -> Aggregate
func ExampleWithoutCommand() {
	id := "123"
	orderCreated := &OrderCreated{
		Model: eventsource.Model{ID: id, Version: 1, At: time.Now()},
	}
	orderApproved := &OrderApproved{
		Model: eventsource.Model{ID: id, Version: 2, At: time.Now()},
	}
	orderShipped := &OrderShipped{
		Model: eventsource.Model{ID: id, Version: 3, At: time.Now()},
	}

	order := Order{}
	order.On(orderCreated)
	order.On(orderApproved)
	order.On(orderShipped)

	fmt.Printf("Order %v on %v\n", order.State, order.UpdatedAt)
}

func ExampleWithDispatcher() {
	serializer := eventsource.NewJSONSerializer(
		OrderCreated{},
		OrderApproved{},
		OrderShipped{},
	)

	store, err := dynamodbstore.New("orders", dynamodbstore.WithRegion("us-west-2"))
	repo := eventsource.New(&Order{},
		eventsource.WithStore(store),
		eventsource.WithSerializer(serializer),
	)

	dispatcher := eventsource.NewDispatcher(repo)
	id := strconv.FormatInt(time.Now().UnixNano(), 36)
	ctx := context.Background()
	order := Order{}

	createOrder := &CreateOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}
	err := dispatcher.Dispatch(ctx, createOrder)
	if err != nil {
		return
	}

	approvedOrder := &ApproveOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}
	err := dispatcher.Dispatch(ctx, approvedOrder)
	if err != nil {
		return
	}

	v, err := repo.Load(ctx, id)
	if err != nil {
		return
	}
	order = v.(*Order)

	fmt.Printf("Order %v on %v\n", order.State, order.UpdatedAt)
}

// Example which create Events out of Command.
func Example() {
	id := "123"
	ctx := context.Background()
	order := Order{}

	createOrder := &CreateOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}
	createEvents, err := order.Apply(ctx, createOrder)
	if err != nil {
		return
	}

	for _, event := range createEvents {
		err := order.On(event)
		if err != nil {
			return
		}
	}

	approvedOrder := &ApproveOrder{
		CommandModel: eventsource.CommandModel{ID: id},
	}
	approvedEvents, err := order.Apply(ctx, approvedOrder)
	if err != nil {
		return
	}

	for _, event := range approvedEvents {
		err := order.On(event)
		if err != nil {
			return
		}
	}

	fmt.Printf("Order %v on %v\n", order.State, order.UpdatedAt)
}

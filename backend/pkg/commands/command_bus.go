package commands

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/mq"
	"golang.org/x/sync/errgroup"
	"strings"
)

type CommandBus interface {
	RegisterHandler(handler CommandHandler) error
}

type RabbitCommandBus struct {
	amqpClient mq.Client
	nameMap    map[string]CommandHandler
	cmdMap     map[string][]CommandHandler
	evtsMap    map[string][]string
	started    bool
	cmdMapper  *CommandMapper
}

func NewRabbitCommandBus(amqpClient mq.Client, commandMapper *CommandMapper) *RabbitCommandBus {
	return &RabbitCommandBus{
		amqpClient: amqpClient,
		nameMap:    map[string]CommandHandler{},
		cmdMap:     map[string][]CommandHandler{},
		evtsMap:    map[string][]string{},
		started:    false,
		cmdMapper:  commandMapper,
	}
}

func (r *RabbitCommandBus) Start(ctx context.Context) error {
	r.started = true
	ch, err := r.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(ctx, "commands.fanout", "fanout", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(ctx, "commands.headers", "headers", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.ExchangeBind(ctx, "commands.headers", "", "commands.fanout", false, nil); err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)

	for handlerName, handler := range r.nameMap {

		g.Go(func() error {

			ch, err := r.amqpClient.GetChannel()
			if err != nil {
				return err
			}

			if err := ch.QueueDeclare(ctx, handlerName, true, false, false, false, map[string]interface{}{}); err != nil {
				return err
			}

			eventsForHandler := r.evtsMap[handlerName]
			for _, commandType := range eventsForHandler {
				if err := ch.QueueBind(ctx, handlerName, "", "commands.headers", false, map[string]interface{}{
					"command_type": commandType,
				}); err != nil {
					return err
				}
			}

			msgChan, err := ch.Consume(ctx, handlerName, "", false, false, false, false, map[string]interface{}{})

			if err != nil {
				return err
			}

			for {
				select {
				case <-ctx.Done():
					return nil
				case msg := <-msgChan:
					cmd, err := r.cmdMapper.Map(msg.Type, msg.Body)
					if err != nil {
						return err
					}
					if err := handler.HandleCommand(ctx, cmd); err != nil {
						return err
					}
					if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
						return err
					}
				}
			}
		})

	}

	return g.Wait()

}

func (r *RabbitCommandBus) RegisterHandler(handler CommandHandler) error {
	if r.started {
		return fmt.Errorf("command bus already started")
	}
	if strings.TrimSpace(handler.GetName()) == "" {
		return fmt.Errorf("handler name is required")
	}
	if strings.TrimSpace(handler.GetName()) != handler.GetName() {
		return fmt.Errorf("invalid handler name")
	}
	if _, ok := r.nameMap[handler.GetName()]; ok {
		return fmt.Errorf("handler with same name already registered: %s", handler.GetName())
	}
	r.nameMap[handler.GetName()] = handler
	// Copy command types arr
	cmdTypes := make([]string, len(handler.GetCommandTypes()))
	copy(cmdTypes, handler.GetCommandTypes())

	// Map handler name to event types
	r.evtsMap[handler.GetName()] = cmdTypes

	for _, cmdType := range cmdTypes {

		// Create slice for handlers that handle event
		if _, ok := r.cmdMap[cmdType]; !ok {
			r.cmdMap[cmdType] = []CommandHandler{}
		}

		// make sure command handler was not already registered
		for _, commandHandler := range r.cmdMap[cmdType] {
			if commandHandler == handler {
				return fmt.Errorf("handler %s already registered for command %s", handler.GetName(), cmdType)
			}
		}

		// append handler to list of handlers for that event
		r.cmdMap[cmdType] = append(r.cmdMap[cmdType], handler)
	}

	return nil
}

var _ CommandBus = &RabbitCommandBus{}

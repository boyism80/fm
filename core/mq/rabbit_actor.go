package mq

import (
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// DirectExchange is the default RabbitMQ direct exchange name used by publishers.
const DirectExchange = "amq.direct"

// RabbitActorConfig is passed when spawning a RabbitActor (one actor per logical consumer).
type RabbitActorConfig struct {
	Root        *actor.RootContext
	Broker      config.RabbitMQEndpoint
	Exchange    string
	QueueName   string
	ConsumerTag string
	RoutingKeys []string
	Dispatcher  *Dispatcher
}

// IncomingAMQP is a copy of the broker payload. The consume goroutine Ack's before Send (at-most-once).
type IncomingAMQP struct {
	Body []byte
}

// RabbitActor owns the AMQP connection, a single consume goroutine, and JSON dispatch on the actor mailbox.
type RabbitActor struct {
	cfg         RabbitActorConfig
	mu          sync.Mutex
	ch          *amqp.Channel
	consumerTag string
	consumeWG   sync.WaitGroup
}

func NewRabbitActor(cfg RabbitActorConfig) *RabbitActor {
	if cfg.Exchange == "" {
		cfg.Exchange = DirectExchange
	}
	return &RabbitActor{cfg: cfg}
}

func (a *RabbitActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		a.consumeWG.Add(1)
		go a.runConsume(ctx.ActorSystem().Root, ctx.Self())
	case *IncomingAMQP:
		if a.cfg.Dispatcher == nil || len(m.Body) == 0 {
			return
		}
		_ = a.cfg.Dispatcher.Dispatch(ctx, m.Body)
	case *actor.Stopping:
		a.shutdownConsume()
	}
}

func (a *RabbitActor) shutdownConsume() {
	a.mu.Lock()
	ch := a.ch
	tag := a.consumerTag
	a.mu.Unlock()
	if ch != nil {
		if err := ch.Cancel(tag, false); err != nil {
			log.Printf("mq: consumer Cancel: %v", err)
		}
	}
	a.consumeWG.Wait()
}

func (a *RabbitActor) runConsume(root *actor.RootContext, self *actor.PID) {
	defer a.consumeWG.Done()
	if root == nil || self == nil {
		log.Printf("mq: RabbitActor runConsume: nil root or self")
		return
	}
	url := a.cfg.Broker.AMQPURL()
	if url == "" {
		log.Printf("mq: RabbitActor: empty AMQP URL, skip consume")
		return
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("mq: Dial: %v", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		log.Printf("mq: Channel: %v", err)
		return
	}
	ex := a.cfg.Exchange
	if ex == DirectExchange {
		if err := ch.ExchangeDeclarePassive(ex, "direct", true, false, false, false, nil); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			log.Printf("mq: ExchangeDeclarePassive: %v", err)
			return
		}
	} else if err := ch.ExchangeDeclare(ex, "direct", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		log.Printf("mq: ExchangeDeclare: %v", err)
		return
	}
	q, err := ch.QueueDeclare(a.cfg.QueueName, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		log.Printf("mq: QueueDeclare: %v", err)
		return
	}
	for _, rk := range a.cfg.RoutingKeys {
		if rk == "" {
			continue
		}
		if err := ch.QueueBind(q.Name, rk, ex, false, nil); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			log.Printf("mq: QueueBind: %v", err)
			return
		}
	}
	deliveries, err := ch.Consume(q.Name, a.cfg.ConsumerTag, false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		log.Printf("mq: Consume: %v", err)
		return
	}

	a.mu.Lock()
	a.ch = ch
	a.consumerTag = a.cfg.ConsumerTag
	a.mu.Unlock()

	defer func() {
		_ = ch.Close()
		_ = conn.Close()
		a.mu.Lock()
		a.ch = nil
		a.mu.Unlock()
	}()

	log.Printf("mq: RabbitActor consuming queue=%s", a.cfg.QueueName)

	for msg := range deliveries {
		if err := msg.Ack(false); err != nil {
			log.Printf("mq: immediate Ack failed: %v", err)
		}
		if len(msg.Body) == 0 {
			continue
		}
		body := make([]byte, len(msg.Body))
		copy(body, msg.Body)
		root.Send(self, &IncomingAMQP{Body: body})
	}
}

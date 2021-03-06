package consumer

import (
	"errors"
	"github.com/streadway/amqp"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	TOPIC_EXCHANGE = "topic"
)

var (
	ERRNAMEREQUIRED = errors.New("name is a required exchange field")
)

// Exchange config sets up a new
// exchange with the provided params
// Defaults are enabled so not all params
// may need set depending on requirements
type ExchangeConfig struct{
	Name string
	exType *string
	durable *bool
	autoDelete *bool
	internal *bool
	args map[string]interface{}
}

type BrokerConfig struct{
	Exchange ExchangeConfig
	Consumers map[string]ConsumerConfig
}

// GetAutoDelete returns the type of deletion policy
// set in config, if nil then it returns a
// default of false
func (e *ExchangeConfig) GetName() (string, error){
	if e.Name != ""{
		return e.Name, nil
	}

	return "", ERRNAMEREQUIRED
}

// GetType returns the type of exchange
// set in config, if nil then it returns a
// default of Topic
func (e *ExchangeConfig) GetType() string{
	if e.exType != nil{
		return *e.exType
	}

	return TOPIC_EXCHANGE
}

// GetDurable returns the type of durability
// set in config, if nil then it returns a
// default of true
func (e *ExchangeConfig) GetDurable() bool{
	if e.durable != nil{
		return *e.durable
	}

	return true
}

// GetAutoDelete returns the type of deletion policy
// set in config, if nil then it returns a
// default of false
func (e *ExchangeConfig) GetAutoDelete() bool{
	if e.autoDelete != nil{
		return *e.autoDelete
	}

	return false
}

// GetInternal determines whether this exchange
// can only be published to from other exchanges
// default value is false meaning external sources
// can by default publish to this exchange
func (e *ExchangeConfig) GetInternal() bool{
	if e.internal != nil{
		return *e.internal
	}

	return false
}

// GetArgs gets a table of arbitrary arguments
// which are passed to the exchange
func (e *ExchangeConfig) GetArgs() map[string]interface{}{
	if e.args == nil{
		return make(map[string]interface{})
	}
	return e.args
}

// BuildExchange builds an exchange
func (e *ExchangeConfig) BuildExchange(ch *amqp.Channel) (err error){
	n, err := e.GetName()
	if err != nil{
		log.Error(err)
		return err
	}

	log.Debugf("setting up %s exchange", n)

	dlx := fmt.Sprintf("%s.deadletter", n)
	if err = ch.ExchangeDeclare(dlx, e.GetType(), e.GetDurable(), e.GetAutoDelete(), e.GetInternal(), false, nil); err != nil{
		log.Errorf("error when setting up deadletter exchange %s: %s", dlx)
		return
	}

	if err = ch.ExchangeDeclare(n, e.GetType(), e.GetDurable(),e.GetAutoDelete(), e.GetInternal(), false, e.GetArgs()); err != nil{
		log.Errorf("error when setting up exchange %s: %s",n, err.Error())
		return
	}
	log.Debugf("%s exchange setup success", n)
	return
}

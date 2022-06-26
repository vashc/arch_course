package prj

const (
	AuthTokenAlgo = "HS256"

	RequestParamUserID  = "user_id"
	RequestParamOrderID = "order_id"
)

const (
	HeaderAuth   = "Authorization"
	HeaderBearer = "Bearer "
)

const (
	CtxAuthToken = "token"
)

const (
	RabbitProtocol    = "amqp"
	RabbitDurable     = true
	RabbitAutoDelete  = false
	RabbitExclusive   = false
	RabbitNoWait      = false
	RabbitNoLocal     = false
	RabbitExchange    = ""
	RabbitMandatory   = false
	RabbitImmediate   = false
	RabbitContentType = "text/plain"
	RabbitConsumer    = ""
	RabbitAutoAck     = true

	QueueSagaSteps       = "saga_steps"
	QueueNotifications   = "notifications"
	QueueExchangeOrders  = "exchange_orders"
	QueueBcgatewayOrders = "bcgateway_orders"
)

const HotWalletUserID = 1

const (
	SagaTypeExchanger = iota
	SagaTypeBcgateway
)

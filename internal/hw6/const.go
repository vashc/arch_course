package hw6

const (
	DBDriver = "postgres"

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

	QueueOrders        = "orders"
	QueueNotifications = "notifications"

	NotificationFailFunds    = "can't complete order processing: insufficient funds"
	NotificationFailInternal = "can't complete order processing: internal error"
	NotificationSuccess      = "complete order processing: success"

	KafkaPartition = 0
)

const (
	AccountID = "account_id"
	Email     = "email"
)

package wallet

const dbDriver = "postgres"

const (
	TypeSellOrder = "sell"
	TypeBuyOrder  = "buy"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

const (
	NotificationOrderCompleted = "Order is successfully processed."
	NotificationOrderFailed    = "Order could not be processed successfully."
)

const CryptoToFiatRatio = 10

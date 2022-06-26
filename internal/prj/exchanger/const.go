package exchanger

const dbDriver = "postgres"

// Order types
const (
	TypeSellOrder = "sell"
	TypeBuyOrder  = "buy"
)

// Exchanger order statuses
const (
	StatusNew       = "new"
	StatusFailed    = "failed"
	StatusCompleted = "completed"
)

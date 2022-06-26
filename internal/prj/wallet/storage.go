package wallet

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq" // Driver
)

func NewStorage(config *Config) (*Storage, error) {
	conn, err := newConn(config)
	if err != nil {
		return nil, err
	}

	sess := conn.NewSession(conn.EventReceiver)

	return &Storage{sess: sess}, nil
}

func newConn(cfg *Config) (*dbr.Connection, error) {
	conn, err := dbr.Open(dbDriver, cfg.DBURI, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errDbrOpenConnection, err.Error())
	}

	return conn, nil
}

func (s *Storage) Close() error {
	return s.sess.Close()
}

func (s *Storage) GetOrderByID(id int64) (order *Order, err error) {
	query := `
SELECT * FROM orders
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.SelectBySql(query, id).LoadOne(&order)
	if err != nil {
		return nil, err
	}

	return order, tx.Commit()
}

func (s *Storage) CreateOrder(order *Order) (orderID int64, err error) {
	query := `
INSERT INTO orders (user_id, type, crypto_amount, fiat_amount, status)
VALUES (?, ?, ?, ?, ?)
RETURNING id;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertBySql(
		query,
		order.UserID,
		order.Type,
		order.CryptoAmount,
		order.FiatAmount,
		order.Status,
	).Load(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, tx.Commit()
}

func (s *Storage) CreateOrderProcessing(order *OrderProcessing) (err error) {
	query := `
INSERT INTO orders_processing (order_id, steps_number)
VALUES (?, ?);
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.InsertBySql(
		query,
		order.OrderID,
		order.StepsNumber,
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) UpdateOrderStatus(orderID int64, status string) error {
	query := `
UPDATE orders
SET status = ?
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.UpdateBySql(query, status, orderID).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DecreaseCurrentStepNumber(orderID int64) (stepsRemain int, err error) {
	query := `
UPDATE orders_processing
SET steps_number = steps_number - 1
WHERE order_id = ?
RETURNING steps_number;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.UpdateBySql(query, orderID).Load(&stepsRemain)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()

	return stepsRemain, err
}

func (s *Storage) FailOrderProcessingStep(orderID int64, step int) (err error) {
	query := `
UPDATE orders_processing
SET failed_steps = array_append(failed_steps, ?)
WHERE order_id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.UpdateBySql(query, step, orderID).Exec()
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (s *Storage) GetFailedOrderProcessingSteps(orderID int64) (steps map[uint8]struct{}, err error) {
	query := `
SELECT UNNEST(failed_steps)
FROM orders_processing
WHERE order_id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var rows *sql.Rows
	rows, err = tx.SelectBySql(query, orderID).Rows()
	if err != nil {
		return nil, err
	}

	tx.SelectBySql(query, orderID)

	var step uint8
	steps = make(map[uint8]struct{})
	for rows.Next() {
		err = rows.Scan(&step)
		if err != nil {
			log.Printf("rows.Scan error: %s", err.Error())
			return nil, err
		}
		steps[step] = struct{}{}
	}

	return steps, tx.Commit()
}

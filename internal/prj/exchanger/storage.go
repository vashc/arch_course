package exchanger

import (
	"arch_course/internal/prj"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"

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

func (s *Storage) CreateOrder(order *prj.ExchangeOrder) (id int64, err error) {
	order.UUID = utils.UUIDv4()

	query := `
INSERT INTO exchange_orders (uuid, acquirer_user_id, order_id, type, crypto_amount, fiat_amount, compensate, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertBySql(
		query,
		order.UUID,
		order.AcquirerUserID,
		order.OrderID,
		order.Type,
		order.CryptoAmount,
		order.FiatAmount,
		order.Compensate,
		order.Status,
	).Load(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (s *Storage) UpdateOrderStatus(id int64, status string) error {
	query := `
UPDATE exchange_orders
SET status = ?	
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.UpdateBySql(query, status, id).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

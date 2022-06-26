package order_service

import (
	"arch_course/internal/hw6"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"time"

	_ "github.com/lib/pq" // Driver
)

func NewStorage(config *hw6.Config) (*Storage, error) {
	conn, err := newConn(config)
	if err != nil {
		return nil, err
	}

	sess := conn.NewSession(conn.EventReceiver)

	return &Storage{Sess: sess}, nil
}

func newConn(cfg *hw6.Config) (*dbr.Connection, error) {
	conn, err := dbr.Open(hw6.DBDriver, cfg.DBURI, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", hw6.ErrDbrOpenConnection, err.Error())
	}

	return conn, nil
}

func (s *Storage) Close() error {
	return s.Sess.Close()
}

func (s *Storage) CreateOrder(order *Order) error {
	query := `
INSERT INTO orders(user_id, price, created_at)
VALUES (?, ?, ?)
RETURNING id;
`

	tx, err := s.Sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertBySql(
		query,
		order.UserID,
		order.Price,
		time.Now(),
	).Load(order)
	if err != nil {
		return err
	}

	return tx.Commit()
}

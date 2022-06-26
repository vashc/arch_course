package notification

import (
	"arch_course/internal/prj"
	"fmt"

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

func (s *Storage) CreateNotification(notification *prj.Notification) (id int64, err error) {
	query := `
INSERT INTO notifications (order_id, email, payload, status)
VALUES (?, ?, ?, ?)
RETURNING id;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertBySql(
		query,
		notification.OrderID,
		notification.Email,
		notification.Payload,
		notification.Status,
	).Load(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (s *Storage) GetNotificationByID(id int64) (notification *prj.Notification, err error) {
	query := `
SELECT *
FROM notifications
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.SelectBySql(query, id).LoadOne(&notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *Storage) UpdateNotificationStatus(id int64, status string) error {
	query := `
UPDATE notifications
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

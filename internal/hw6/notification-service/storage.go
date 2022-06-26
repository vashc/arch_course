package notification_service

import (
	"fmt"
	"time"

	"arch_course/internal/hw6"

	"github.com/gocraft/dbr/v2"
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

func (s *Storage) SendMessage(message *Message) error {
	query := `
INSERT INTO messages(email, payload, created_at)
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
		message.Email,
		message.Payload,
		time.Now(),
	).Load(message)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetMessagesByEmail(email string) (message []Message, err error) {
	query := `
SELECT *
FROM messages
WHERE email = ?;
`

	tx, err := s.Sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	var messages []Message
	_, err = tx.SelectBySql(query, email).Load(&messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

package auth

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

func (s *Storage) RegisterUser(user *prj.User) error {
	query := `
INSERT INTO users(username, first_name, last_name, email, password)
VALUES (?, ?, ?, ?, ?)
RETURNING id;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertBySql(
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Load(user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetUserByUsername(username, password string) (user *prj.User, err error) {
	query := `
SELECT *
FROM users
WHERE username = ? AND password = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.SelectBySql(query, username, password).LoadOne(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) GetUserByID(id int64) (user *prj.User, err error) {
	query := `
SELECT *
FROM users
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.SelectBySql(query, id).LoadOne(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

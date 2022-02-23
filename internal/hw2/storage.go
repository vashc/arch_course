package hw2

import (
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

func (s *Storage) CreateUser(user *User) error {
	query := `
INSERT INTO users(username, first_name, last_name, email, phone)
VALUES (?, ?, ?, ?, ?);
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.InsertBySql(
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetUser(id int64) (user *User, err error) {
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

func (s *Storage) UpdateUser(user *User) error {
	// It's better to create a query constructor in here
	query := `
UPDATE users u
SET 
	username = COALESCE(NULLIF(?, ''), u.username),
	first_name = COALESCE(NULLIF(?, ''), u.first_name),	
	last_name = COALESCE(NULLIF(?, ''), u.last_name),	
	email = COALESCE(NULLIF(?, ''), u.email),	
	phone = COALESCE(NULLIF(?, ''), u.phone)
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.UpdateBySql(
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.ID,
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DeleteUser(id int64) error {
	query := `
DELETE FROM users
WHERE id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.DeleteBySql(query, id).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

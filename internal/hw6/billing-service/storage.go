package billing_service

import (
	"arch_course/internal/hw6"
	"database/sql"
	"fmt"
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

func (s *Storage) CreateAccount(account *Account) error {
	query := `
INSERT INTO accounts(username, email, password, balance)
VALUES (?, ?, ?, ?)
RETURNING id;
`

	tx, err := s.Sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.InsertBySql(
		query,
		account.Username,
		account.Email,
		account.Password,
		account.Balance,
	).Load(account)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetAccountByID(id int64) (account *Account, err error) {
	query := `
SELECT *
FROM accounts
WHERE id = ?;
`

	tx, err := s.Sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.SelectBySql(query, id).LoadOne(&account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Storage) UpdateAccountBalance(account *Account) (err error) {
	query := `
UPDATE accounts
SET balance = ?
WHERE id = ?;
`

	tx, err := s.Sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.UpdateBySql(query, account.Balance, account.ID).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) CompareAndUpdateAccountBalance(account *Account, balance float64) (err error) {
	query := `
UPDATE accounts a
SET balance = a.balance - ?
WHERE id = ? AND a.balance >= ?;
`

	tx, err := s.Sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	var res sql.Result
	res, err = tx.UpdateBySql(query, balance, account.ID, balance).Exec()
	if err != nil {
		return err
	}

	var updated int64
	updated, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return hw6.ErrInsufficientFunds
	}

	return tx.Commit()
}

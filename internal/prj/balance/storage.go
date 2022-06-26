package balance

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

func (s *Storage) CreateWallet(wallet *prj.Wallet) (err error) {
	query := `
INSERT INTO wallets (user_id, crypto_amount, fiat_amount)
VALUES (?, ?, ?);
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.InsertBySql(
		query,
		wallet.UserID,
		wallet.CryptoAmount,
		wallet.FiatAmount,
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetWalletByUserID(userID int64) (wallet *prj.Wallet, err error) {
	query := `
SELECT *
FROM wallets
WHERE user_id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()

	err = tx.SelectBySql(query, userID).LoadOne(&wallet)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *Storage) UpdateWallet(wallet *prj.Wallet) error {
	// It's better to create a query constructor in here
	query := `
UPDATE wallets
SET 
	crypto_amount = ?,
	fiat_amount = ?
WHERE user_id = ?;
`

	tx, err := s.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.UpdateBySql(
		query,
		wallet.CryptoAmount,
		wallet.FiatAmount,
		wallet.UserID,
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

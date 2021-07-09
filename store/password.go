package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/faceit/test/entity"
)

// user password parameters and query
const (
	passwordTable  = `users_password`
	passwordParams = `password_id, pwd, salt`

	createPasswordQuery = `INSERT INTO ` + passwordTable + ` ( ` + passwordParams + ` ) VALUES ($1, $2, $3);`
	selectPassqordQuery = `SELECT ` + passwordParams + ` FROM ` + passwordTable + ` WHERE password_id = $1;`
	updatePasswordQuery = `UPDATE ` + passwordTable + ` SET pwd = $1 AND salt = $2 WHERE password_id = $3;`
	deletePassqordQuery = `DELETE FROM ` + passwordTable + ` WHERE password_id = $1;`
)

// Password is a pasword store implementation
type Password struct {
	*sql.DB
}

// NewPassword creates a new password instance
func NewPassword(db *sql.DB) *Password {
	return &Password{
		db,
	}
}

// Update updates users_password record in database by id
func (p *Password) Update(ctx context.Context, userID int, hash, salt string) error {
	_, err := p.ExecContext(ctx, updatePasswordQuery, hash, salt, userID)
	if err != nil {
		return fmt.Errorf("query failed, %w", err)
	}

	return nil
}

// One returns one record from users_password by id
func (p *Password) One(ctx context.Context, id int) (entity.Password, error) {
	pwd := entity.Password{}

	err := p.QueryRowContext(ctx, selectPassqordQuery, id).Scan(
		&pwd.UserID,
		&pwd.Hash,
		&pwd.Salt)
	if err != nil {
		err = fmt.Errorf("query failed, %w", err)
	}

	return pwd, err
}

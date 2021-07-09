package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/faceit/test/entity"
)

// user table parameters and query
const (
	userTable = `users`

	createUserQuery = `INSERT INTO ` +
		userTable + ` (first_name, last_name, nick_name, email, country) VALUES ($1, $2, $3, $4, $5) RETURNING user_id;`

	updateUserQuery = `UPDATE ` +
		userTable + ` SET first_name = $1 AND last_name = $2 AND nick_name = $3 AND email = $4 WHERE user_id = $5;`

	deleteUserQuery = `DELETE FROM ` + userTable + ` WHERE user_id = $1;`

	selectOneUserQuery = `SELECT u.user_id, u.first_name, u.last_name, u.nick_name, u.email, c.country_name FROM ` +
		userTable + ` as u, ` + countryTable + ` as c WHERE u.user_id = $1 AND c.country_id = u.country;`

	selectAllUsersQuery = `SELECT u.user_id, u.first_name, u.last_name, u.nick_name, u.email, c.country_name FROM ` +
		userTable + ` as u, ` + countryTable + ` as c WHERE c.country_id = u.country;`

	selectAllUsersByCountryQuery = `SELECT u.user_id, u.first_name, u.last_name, u.nick_name, u.email, c.country_name FROM ` +
		userTable + ` as u, ` + countryTable + ` as c WHERE c.country_name = $1 AND c.country_id = u.country;`

	selectAllUsersByFilterQuery = `SELECT u.user_id, u.first_name, u.last_name, u.nick_name, u.email, c.country_name FROM ` +
		userTable + ` as u, ` + countryTable + ` as c WHERE u.%s = $1 AND c.country_id = u.country;`
)

// User is a user store implementation
type User struct {
	*sql.DB
	tx *sql.TxOptions
}

// NewUser creates a new User instance
func NewUser(db *sql.DB, tx *sql.TxOptions) *User {
	return &User{
		db,
		tx,
	}
}

// Create creates a new users record in database
func (u *User) Create(ctx context.Context, user entity.User) (int, error) {
	var id int

	// starting a db transaction
	tx, err := u.BeginTx(ctx, u.tx)
	if err != nil {
		return id, fmt.Errorf("failed to begin transaction, %w", err)
	}

	// creating user and parsing user_id into id var for furthure password creation
	err = tx.QueryRowContext(ctx, createUserQuery, user.FirstName, user.LastName, user.NickName, user.Email, user.CountryID).Scan(&id)
	if err != nil {
		return 0, u.rollbackTransaction(tx, err)
	}

	// creating password for user with id from last query
	_, err = tx.ExecContext(ctx, createPasswordQuery, id, user.Password, user.Salt)
	if err != nil {
		return 0, u.rollbackTransaction(tx, err)
	}

	// commititng TX
	err = tx.Commit()
	if err != nil {
		return 0, u.rollbackTransaction(tx, err)
	}

	return id, err
}

// Update Updates a users record in database by it's id
func (u *User) Update(ctx context.Context, user entity.User) error {
	_, err := u.ExecContext(ctx, updateUserQuery, user.FirstName, user.LastName, user.NickName, user.Email, user.CountryID)
	if err != nil {
		return fmt.Errorf("query failed, %w", err)
	}

	return nil
}

// Delete deletes a users record from database by it's id
func (u *User) Delete(ctx context.Context, id int) error {
	// starting a db transaction
	tx, err := u.BeginTx(ctx, u.tx)
	if err != nil {
		return fmt.Errorf("begin transaction failed, %w", err)
	}

	// deleting user by his id
	_, err = tx.ExecContext(ctx, deletePassqordQuery, id)
	if err != nil {
		return u.rollbackTransaction(tx, err)
	}

	// deleting user's password
	_, err = tx.ExecContext(ctx, deleteUserQuery, id)
	if err != nil {
		return u.rollbackTransaction(tx, err)
	}

	// commititng TX
	err = tx.Commit()
	if err != nil {
		return u.rollbackTransaction(tx, err)
	}

	return nil
}

// One returns one users record from database by id
func (u *User) One(ctx context.Context, id int) (entity.User, error) {
	user := entity.User{}

	err := u.QueryRowContext(ctx, selectOneUserQuery, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.NickName,
		&user.Email,
		&user.Country)
	if errors.Is(err, sql.ErrNoRows) {
		return user, entity.ErrNotFound
	}
	if err != nil {
		return user, fmt.Errorf("query failed, %w", err)
	}

	return user, err
}

// All returns all users records from database
func (u *User) All(ctx context.Context) ([]entity.User, error) {
	userRows, err := u.QueryContext(ctx, selectAllUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("query failed, %w", err)
	}

	defer func() {
		_ = userRows.Close()
	}()

	users := []entity.User{}

	for userRows.Next() {
		user := entity.User{}

		err = userRows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.NickName,
			&user.Email,
			&user.Country)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("scan results failed, %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

// AllByCountry gets all users for selected country
func (u *User) AllByCountry(ctx context.Context, iso2 string) ([]entity.User, error) {
	userRows, err := u.QueryContext(ctx, selectAllUsersByCountryQuery, iso2)
	if err != nil {
		return nil, fmt.Errorf("query failed, %w", err)
	}

	defer func() {
		_ = userRows.Close()
	}()

	users := []entity.User{}

	for userRows.Next() {
		user := entity.User{}

		err = userRows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.NickName,
			&user.Email,
			&user.Country)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("query failed, %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

// AllWithFilter gets all users by selected filter
func (u *User) AllWithFilter(ctx context.Context, title, filter string) ([]entity.User, error) {
	userRows, err := u.QueryContext(ctx, fmt.Sprintf(selectAllUsersByFilterQuery, title), filter)
	if err != nil {
		return nil, fmt.Errorf("query failed, %w", err)
	}

	defer func() {
		_ = userRows.Close()
	}()

	users := []entity.User{}

	for userRows.Next() {
		user := entity.User{}

		err = userRows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.NickName,
			&user.Email,
			&user.Country)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("query failed, %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *User) rollbackTransaction(tx *sql.Tx, e error) error {
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("%w, rollback transaction faileu. error: %s", e, err)
	}

	return e
}

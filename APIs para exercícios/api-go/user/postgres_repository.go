package user

import (
	"database/sql"
	"errors"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func MigratePostgres(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE
);
`)
	return err
}

func (r *PostgresUserRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query(`SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *PostgresUserRepository) GetByID(id int) (*User, error) {
	var u User
	err := r.db.QueryRow(`SELECT id, name, email FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Create(u User) (*User, error) {
	err := r.db.QueryRow(
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		u.Name, u.Email,
	).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Update(id int, u User) (*User, error) {
	res, err := r.db.Exec(`UPDATE users SET name = $1, email = $2 WHERE id = $3`, u.Name, u.Email, id)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("user not found")
	}
	u.ID = id
	return &u, nil
}

func (r *PostgresUserRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("user not found")
	}
	return nil
}

package repository

import (
    "database/sql"
    "fmt"

    "github.com/pkgzx/liliApi/src/pkg/data"
)

type UserRepository struct {
    *BaseRepository
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *UserRepository) GetByUsername(username string) (*data.User, error) {
    query := `
        SELECT id, username, password_hash, full_name, created_at
        FROM users 
        WHERE username = $1
    `
    
    var user data.User
    err := r.db.QueryRow(query, username).Scan(
        &user.ID,
        &user.Username,
        &user.PasswordHash,
        &user.FullName,
        &user.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("error getting user: %w", err)
    }

    return &user, nil
}

func (r *UserRepository) Create(user *data.User) error {
    query := `
        INSERT INTO users (username, password_hash, full_name)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
    
    err := r.db.QueryRow(
        query,
        user.Username,
        user.PasswordHash,
        user.FullName,
    ).Scan(&user.ID, &user.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("error creating user: %w", err)
    }

    return nil
}

func (r *UserRepository) Update(user *data.User) error {
    query := `
        UPDATE users 
        SET username = $2, full_name = $3
        WHERE id = $1
    `
    
    result, err := r.db.Exec(query, user.ID, user.Username, user.FullName)
    if err != nil {
        return fmt.Errorf("error updating user: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

func (r *UserRepository) UpdatePassword(id int32, passwordHash string) error {
    query := `
        UPDATE users 
        SET password_hash = $2
        WHERE id = $1
    `
    
    result, err := r.db.Exec(query, id, passwordHash)
    if err != nil {
        return fmt.Errorf("error updating password: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}
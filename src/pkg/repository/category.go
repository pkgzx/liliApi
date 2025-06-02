package repository

import (
    "database/sql"
    "fmt"

    "github.com/pkgzx/liliApi/src/pkg/data"
)

type CategoryRepository struct {
    *BaseRepository
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
    return &CategoryRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *CategoryRepository) GetAll() ([]data.Category, error) {
    query := `
        SELECT id, name, created_at 
        FROM categories 
        ORDER BY name
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error querying categories: %w", err)
    }
    defer rows.Close()

    var categories []data.Category
    if err := ScanRowsToStruct(rows, &categories); err != nil {
        return nil, fmt.Errorf("error scanning categories: %w", err)
    }

    return categories, nil
}

func (r *CategoryRepository) GetByID(id int32) (*data.Category, error) {
    query := `
        SELECT id, name, created_at 
        FROM categories 
        WHERE id = $1
    `
    
    var category data.Category
    err := r.db.QueryRow(query, id).Scan(
        &category.ID,
        &category.Name,
        &category.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("error getting category: %w", err)
    }

    return &category, nil
}
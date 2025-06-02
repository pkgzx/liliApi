package repository

import (
    "database/sql"
    "fmt"

    "github.com/pkgzx/liliApi/src/pkg/data"
)

type ProductRepository struct {
    *BaseRepository
}

func NewProductRepository(db *sql.DB) *ProductRepository {
    return &ProductRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *ProductRepository) GetAll() ([]data.Product, error) {
    query := `
        SELECT id, name, description, price, category_id, image_url, is_available, created_at 
        FROM products 
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error querying products: %w", err)
    }
    defer rows.Close()

    var products []data.Product
    if err := ScanRowsToStruct(rows, &products); err != nil {
        return nil, fmt.Errorf("error scanning products: %w", err)
    }

    return products, nil
}

func (r *ProductRepository) GetByID(id int32) (*data.Product, error) {
    query := `
        SELECT id, name, description, price, category_id, image_url, is_available, created_at 
        FROM products 
        WHERE id = $1
    `
    
    var product data.Product
    err := r.db.QueryRow(query, id).Scan(
        &product.ID,
        &product.Name,
        &product.Description,
        &product.Price,
        &product.CategoryID,
        &product.ImageURL,
        &product.IsAvailable,
        &product.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("error getting product: %w", err)
    }

    return &product, nil
}

func (r *ProductRepository) GetByCategory(categoryID int32) ([]data.Product, error) {
    query := `
        SELECT id, name, description, price, category_id, image_url, is_available, created_at 
        FROM products 
        WHERE category_id = $1 AND is_available = true
        ORDER BY name
    `
    
    rows, err := r.db.Query(query, categoryID)
    if err != nil {
        return nil, fmt.Errorf("error querying products by category: %w", err)
    }
    defer rows.Close()

    var products []data.Product
    if err := ScanRowsToStruct(rows, &products); err != nil {
        return nil, fmt.Errorf("error scanning products: %w", err)
    }

    return products, nil
}

func (r *ProductRepository) Create(product *data.Product) error {
    query := `
        INSERT INTO products (name, description, price, category_id, image_url, is_available)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
    `
    
    err := r.db.QueryRow(
        query,
        product.Name,
        product.Description,
        product.Price,
        product.CategoryID,
        product.ImageURL,
        product.IsAvailable,
    ).Scan(&product.ID, &product.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("error creating product: %w", err)
    }

    return nil
}

func (r *ProductRepository) Update(product *data.Product) error {
    query := `
        UPDATE products 
        SET name = $2, description = $3, price = $4, category_id = $5, 
            image_url = $6, is_available = $7
        WHERE id = $1
    `
    
    result, err := r.db.Exec(
        query,
        product.ID,
        product.Name,
        product.Description,
        product.Price,
        product.CategoryID,
        product.ImageURL,
        product.IsAvailable,
    )
    
    if err != nil {
        return fmt.Errorf("error updating product: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("product not found")
    }

    return nil
}

func (r *ProductRepository) Delete(id int32) error {
    query := "DELETE FROM products WHERE id = $1"
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("error deleting product: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("product not found")
    }

    return nil
}
package repository

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/pkgzx/liliApi/src/pkg/data"
)

type OrderRepository struct {
    *BaseRepository
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
    return &OrderRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

func (r *OrderRepository) GetAll() ([]data.Order, error) {
    query := `
        SELECT id, order_number, status, total_amount, notes, created_at, updated_at
        FROM orders 
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error querying orders: %w", err)
    }
    defer rows.Close()

    var orders []data.Order
    if err := ScanRowsToStruct(rows, &orders); err != nil {
        return nil, fmt.Errorf("error scanning orders: %w", err)
    }

    return orders, nil
}

func (r *OrderRepository) GetByID(id int32) (*data.Order, error) {
    query := `
        SELECT id, order_number, status, total_amount, notes, created_at, updated_at
        FROM orders 
        WHERE id = $1
    `
    
    var order data.Order
    err := r.db.QueryRow(query, id).Scan(
        &order.ID,
        &order.OrderNumber,
        &order.Status,
        &order.TotalAmount,
        &order.Notes,
        &order.CreatedAt,
        &order.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("error getting order: %w", err)
    }

    return &order, nil
}

func (r *OrderRepository) Create(order *data.Order) error {
    tx, err := r.db.Begin()
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback()

    // Generar número de orden único
    orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())
    
    query := `
        INSERT INTO orders (order_number, status, total_amount, notes)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
    
    err = tx.QueryRow(
        query,
        orderNumber,
        order.Status,
        order.TotalAmount,
        order.Notes,
    ).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("error creating order: %w", err)
    }

    order.OrderNumber = orderNumber

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("error committing transaction: %w", err)
    }

    return nil
}

func (r *OrderRepository) UpdateStatus(id int32, status string) error {
    query := `
        UPDATE orders 
        SET status = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `
    
    result, err := r.db.Exec(query, id, status)
    if err != nil {
        return fmt.Errorf("error updating order status: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("order not found")
    }

    return nil
}

func (r *OrderRepository) GetByStatus(status string) ([]data.Order, error) {
    query := `
        SELECT id, order_number, status, total_amount, notes, created_at, updated_at
        FROM orders 
        WHERE status = $1
        ORDER BY created_at ASC
    `
    
    rows, err := r.db.Query(query, status)
    if err != nil {
        return nil, fmt.Errorf("error querying orders by status: %w", err)
    }
    defer rows.Close()

    var orders []data.Order
    if err := ScanRowsToStruct(rows, &orders); err != nil {
        return nil, fmt.Errorf("error scanning orders: %w", err)
    }

    return orders, nil
}
package data

import "time"

type User struct {
    ID          int32     `json:"id" db:"id"`
    Username    string    `json:"username" db:"username"`
    PasswordHash string   `json:"password_hash" db:"password_hash"`
    FullName    string    `json:"full_name" db:"full_name"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Category struct {
    ID        int32     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Product struct {
    ID          int32     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    Price       float64   `json:"price" db:"price"`
    CategoryID  int32     `json:"category_id" db:"category_id"`
    ImageURL    string    `json:"image_url" db:"image_url"`
    IsAvailable bool      `json:"is_available" db:"is_available"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Ingredient struct {
    ID            int32     `json:"id" db:"id"`
    Name          string    `json:"name" db:"name"`
    Unit          string    `json:"unit" db:"unit"`
    StockQuantity float64   `json:"stock_quantity" db:"stock_quantity"`
    MinStock      float64   `json:"min_stock" db:"min_stock"`
    CostPerUnit   float64   `json:"cost_per_unit" db:"cost_per_unit"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Order struct {
    ID          int32     `json:"id" db:"id"`
    OrderNumber string    `json:"order_number" db:"order_number"`
    Status      string    `json:"status" db:"status"`
    TotalAmount float64   `json:"total_amount" db:"total_amount"`
    Notes       string    `json:"notes" db:"notes"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
    ID        int32   `json:"id" db:"id"`
    OrderID   int32   `json:"order_id" db:"order_id"`
    ProductID int32   `json:"product_id" db:"product_id"`
    Quantity  int32   `json:"quantity" db:"quantity"`
    UnitPrice float64 `json:"unit_price" db:"unit_price"`
    Subtotal  float64 `json:"subtotal" db:"subtotal"`
}

type InventoryMovement struct {
    ID           int32     `json:"id" db:"id"`
    IngredientID int32     `json:"ingredient_id" db:"ingredient_id"`
    MovementType string    `json:"movement_type" db:"movement_type"` // "entrada", "salida", "ajuste"
    Quantity     float64   `json:"quantity" db:"quantity"`
    Reason       string    `json:"reason" db:"reason"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
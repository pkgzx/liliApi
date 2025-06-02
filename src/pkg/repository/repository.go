package repository

import (
    "database/sql"
    "fmt"
    "reflect"
    "strings"
)

type BaseRepository struct {
    db *sql.DB
}

func NewBaseRepository(db *sql.DB) *BaseRepository {
    return &BaseRepository{db: db}
}

func (r *BaseRepository) GetDB() *sql.DB {
    return r.db
}

// Función auxiliar para escanear filas a structs
func ScanRowsToStruct(rows *sql.Rows, dest any) error {
    v := reflect.ValueOf(dest).Elem()
    if v.Kind() != reflect.Slice {
        return fmt.Errorf("dest must be a pointer to a slice")
    }

    elemType := v.Type().Elem()
    if elemType.Kind() == reflect.Ptr {
        elemType = elemType.Elem()
    }

    columns, err := rows.Columns()
    if err != nil {
        return err
    }

    for rows.Next() {
        elem := reflect.New(elemType).Elem()
        values := make([]any, len(columns))
        
        for i, col := range columns {
            field := findFieldByDBTag(elem, col)
            if field.IsValid() {
                values[i] = field.Addr().Interface()
            } else {
                var dummy interface{}
                values[i] = &dummy
            }
        }

        if err := rows.Scan(values...); err != nil {
            return err
        }

        if elemType.Kind() == reflect.Ptr {
            v.Set(reflect.Append(v, elem.Addr()))
        } else {
            v.Set(reflect.Append(v, elem))
        }
    }

    return rows.Err()
}

func findFieldByDBTag(v reflect.Value, tagValue string) reflect.Value {
    t := v.Type()
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        if tag := field.Tag.Get("db"); tag == tagValue {
            return v.Field(i)
        }
    }
    return reflect.Value{}
}

// Construir query SELECT dinámicamente
func BuildSelectQuery(tableName string, conditions map[string]interface{}) (string, []interface{}) {
    query := fmt.Sprintf("SELECT * FROM %s", tableName)
    args := make([]any, 0)
    
    if len(conditions) > 0 {
        whereClauses := make([]string, 0)
        argIndex := 1
        
        for field, value := range conditions {
            whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", field, argIndex))
            args = append(args, value)
            argIndex++
        }
        
        query += " WHERE " + strings.Join(whereClauses, " AND ")
    }
    
    return query, args
}
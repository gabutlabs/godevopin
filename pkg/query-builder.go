package pkg

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type FilterParams struct {
	Filters map[string]map[string]string `query:"filters"`
}

type QueryParams struct {
	FilterParams
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

type QueryBuilder struct {
	Filters        map[string]map[string]string
	Query          *gorm.DB
	AllowedColumns map[string]string
}

func NewQueryParams(filters map[string]map[string]string, query *gorm.DB, allowedColumns map[string]string) *QueryBuilder {
	return &QueryBuilder{
		Filters:        filters,
		Query:          query,
		AllowedColumns: allowedColumns,
	}
}

// ApplyFilters adalah fungsi global untuk menambahkan filter dinamis ke query GORM.
// Ini aman dari SQL Injection karena menggunakan whitelist 'allowedColumns'.
func (instance *QueryBuilder) ApplyFilters() *gorm.DB {

	// Iterasi semua filter yang diminta dari URL
	for columnKey, conditions := range instance.Filters {

		// 1. KEAMANAN: Validasi kolom dengan whitelist
		// 'columnKey' adalah dari URL (e.g., "name")
		// 'dbColumn' adalah nama kolom di DB (e.g., "user_name")
		dbColumn, isAllowed := instance.AllowedColumns[columnKey]
		if !isAllowed {
			// Jika kolom tidak diizinkan, abaikan filter ini
			continue
		}

		// 2. Iterasi semua kondisi (eq, gte, like, dll)
		for condition, value := range conditions {
			var queryStr string
			var queryVal interface{} = value

			switch condition {
			case "eq":
				queryStr = fmt.Sprintf("%s = ?", dbColumn)
			case "neq":
				queryStr = fmt.Sprintf("%s != ?", dbColumn)
			case "gt":
				queryStr = fmt.Sprintf("%s > ?", dbColumn)
			case "gte":
				queryStr = fmt.Sprintf("%s >= ?", dbColumn)
			case "lt":
				queryStr = fmt.Sprintf("%s < ?", dbColumn)
			case "lte":
				queryStr = fmt.Sprintf("%s <= ?", dbColumn)
			case "like":
				queryStr = fmt.Sprintf("%s LIKE ?", dbColumn)
				queryVal = "%" + value + "%"
			case "in":
				queryStr = fmt.Sprintf("%s IN (?)", dbColumn)
				queryVal = strings.Split(value, ",") // GORM butuh slice untuk IN
			default:
				// Abaikan kondisi yang tidak dikenal
				continue
			}

			// 3. Tambahkan ke query GORM
			// Ini aman karena 'dbColumn' berasal dari whitelist, bukan input user
			instance.Query = instance.Query.Where(queryStr, queryVal)
		}
	}

	// 4. Kembalikan query yang sudah dimodifikasi
	return instance.Query
}

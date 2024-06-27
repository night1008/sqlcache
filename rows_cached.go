package sqlcache

import (
	"database/sql/driver"
	"io"

	"github.com/prashanthpai/sqlcache/cache"
)

// RowsCached implements driver.Rows interface
type RowsCached struct {
	*cache.Item
	ptr int
}

func (r *RowsCached) Columns() []string {
	return r.Item.Cols
}

func (r *RowsCached) ColumnTypeDatabaseTypeName(index int) string {
	if index < len(r.Item.DatabaseTypeNames) {
		return r.Item.DatabaseTypeNames[index]
	}
	return ""
}

func (r *RowsCached) Next(dest []driver.Value) error {
	if r.ptr >= len(r.Item.Rows) {
		return io.EOF
	}

	for i := range dest {
		dest[i] = r.Item.Rows[r.ptr][i]
	}
	r.ptr++

	return nil
}

func (r *RowsCached) Close() error {
	return nil
}

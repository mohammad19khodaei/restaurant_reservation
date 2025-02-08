package table

import "context"

type Repository interface {
	CreateTable(ctx context.Context, table *Table) error
	GetTotalCount(ctx context.Context) int
	CreateTableSettings(ctx context.Context, seatPrice int) error
}

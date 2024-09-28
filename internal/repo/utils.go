package repo

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func updateRow(db *sqlx.DB, table, pk_column string, pk any, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	sql, args := squirrel.Update(table).
		SetMap(fields).
		Where(squirrel.Eq{pk_column: pk}).
		PlaceholderFormat(squirrel.Dollar).
		MustSql()

	result, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}

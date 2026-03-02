package mapper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/walmaa/skemr-api/db/sqlc"
	"github.com/walmaa/skemr-api/internal/dto"
)

func Text(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{
			String: "",
			Valid:  false,
		}
	}
	return pgtype.Text{
		String: *v,
		Valid:  true,
	}
}

func TextPtr(v *pgtype.Text) *string {
	if v != nil && v.Valid {
		return &v.String
	}
	return nil
}

func Int4(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{
			Int32: 0,
			Valid: false,
		}
	}
	return pgtype.Int4{
		Int32: *v,
		Valid: true,
	}
}

func Time(v *pgtype.Timestamptz) time.Time {
	if v.Valid {
		return v.Time
	}
	return time.Time{}
}

func TimePtr(v *pgtype.Timestamptz) *time.Time {
	if v != nil && v.Valid {
		return &v.Time
	}
	return nil
}

func NullDatabaseType(databaseType dto.DatabaseType) sqlc.NullDatabaseType {
	if databaseType == "" {
		return sqlc.NullDatabaseType{
			DatabaseType: "",
			Valid:        false,
		}
	}
	return sqlc.NullDatabaseType{
		DatabaseType: sqlc.DatabaseType(databaseType),
		Valid:        true,
	}
}

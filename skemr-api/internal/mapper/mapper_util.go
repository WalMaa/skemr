package mapper

import "github.com/jackc/pgx/v5/pgtype"

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

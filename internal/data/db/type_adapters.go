package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func IntoPgTimePrt(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{
			Time:             time.Time{},
			InfinityModifier: 0,
			Valid:            false,
		}
	}

	return pgtype.Timestamp{
		Time:             *t,
		InfinityModifier: 0,
		Valid:            true,
	}
}

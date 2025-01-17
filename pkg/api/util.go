package api

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}

	return &ns.String
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}

	return &nt.Time
}

func nullUUIDToUUIDPtr(nu uuid.NullUUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.UUID
}

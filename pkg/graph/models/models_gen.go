// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
)

type BigQuerySource struct {
	// table is the name of the BigQuery table.
	Table string `json:"table"`
	// dataset is the name of the BigQuery dataset.
	Dataset string `json:"dataset"`
}

type AccessRequestStatus string

const (
	AccessRequestStatusPending  AccessRequestStatus = "pending"
	AccessRequestStatusApproved AccessRequestStatus = "approved"
	AccessRequestStatusDenied   AccessRequestStatus = "denied"
)

var AllAccessRequestStatus = []AccessRequestStatus{
	AccessRequestStatusPending,
	AccessRequestStatusApproved,
	AccessRequestStatusDenied,
}

func (e AccessRequestStatus) IsValid() bool {
	switch e {
	case AccessRequestStatusPending, AccessRequestStatusApproved, AccessRequestStatusDenied:
		return true
	}
	return false
}

func (e AccessRequestStatus) String() string {
	return string(e)
}

func (e *AccessRequestStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccessRequestStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccessRequestStatus", str)
	}
	return nil
}

func (e AccessRequestStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

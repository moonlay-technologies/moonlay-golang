package models

import (
	"order-service/global/utils/model"
	"time"
)

type SalesmanAssignment struct {
	ID             int        `json:"id,omitempty"`
	MappingStoreID int        `json:"mapping_store_id,omitempty"`
	SalesmanID     int        `json:"salesman_id,omitempty"`
	BrandID        int        `json:"brand_id,omitempty"`
	AgentID        int        `json:"agent_id,omitempty"`
	DateAssign     *time.Time `json:"date_assign,omitempty"`
	EventSource    string     `json:"event_source,omitempty"`
	EventTransfer  string     `json:"event_transfer,omitempty"`
}

type SalesmanAssignmentChan struct {
	SalesmanAssignment *SalesmanAssignment
	Total              int64
	Error              error
	ErrorLog           *model.ErrorLog
	ID                 int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanAssignmentsChan struct {
	SalesmanAssignments []*SalesmanAssignment
	Total               int64
	Error               error
	ErrorLog            *model.ErrorLog
	ID                  int64 `json:"id,omitempty" bson:"id,omitempty"`
}

type SalesmanAssignmentRequest struct {
	PerPage int    `json:"per_page,omitempty" bson:"per_page,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}

type SalesmanAssignments struct {
	SalesmanAssignments []*SalesmanAssignment `json:"salesman_assignments,omitempty"`
	Total               int64                 `json:"total,omitempty"`
}

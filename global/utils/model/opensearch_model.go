package model

type OpenSearchInsertResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int    `json:"_version,omitempty"`
	Result  string `json:"result,omitempty"`
}

type OpenSearchGetResponse struct {
	Index       string      `json:"_index,omitempty"`
	Type        string      `json:"_type,omitempty"`
	ID          string      `json:"_id,omitempty"`
	Version     int         `json:"_version,omitempty"`
	SecNo       int         `json:"_sec_no,omitempty"`
	PrimaryTerm int         `json:"_primary_term,omitempty"`
	Found       bool        `json:"found,omitempty"`
	Source      interface{} `json:"_source,omitempty"`
}

type OpenSearchCountResponse struct {
	Count int64 `json:"count,omitempty"`
}

type OpenSearchQueryResponse struct {
	Took    int                     `json:"took,omitempty"`
	TimeOut bool                    `json:"time_out,omitempty"`
	Hits    *OpenSearchHitsResponse `json:"hits,omitempty"`
}

type OpenSearchHitsResponse struct {
	Total    OpenSearchHitsTotal           `json:"total,omitempty"`
	MaxScore float64                       `json:"max_score,omitempty"`
	Hits     []*OpenSearchHitsItemResponse `json:"hits,omitempty"`
}

type OpenSearchHitsItemResponse struct {
	Index  string      `json:"_index,omitempty"`
	ID     string      `json:"_id,omitempty"`
	Score  float64     `json:"_score,omitempty"`
	Source interface{} `json:"_source,omitempty"`
}

type OpenSearchHitsTotal struct {
	Value    int    `json:"value,omitempty"`
	Relation string `json:"relation,omitempty"`
}

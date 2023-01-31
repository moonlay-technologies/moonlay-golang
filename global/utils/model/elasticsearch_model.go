package model

type ElasticSearchInsertResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int    `json:"_version,omitempty"`
	Result  string `json:"result,omitempty"`
}

type ElasticSearchQueryResponse struct {
	Took    float64                    `json:"took,omitempty"`
	TimeOut bool                       `json:"time_out,omitempty"`
	Hits    *ElasticSearchHitsResponse `json:"hits,omitempty"`
}

type ElasticSearchHitsResponse struct {
	Total    ElasticSearchHitsTotal           `json:"total,omitempty"`
	MaxScore float64                          `json:"max_score,omitempty"`
	Hits     []*ElasticSearchHitsItemResponse `json:"hits,omitempty"`
}

type ElasticSearchHitsItemResponse struct {
	Index  string      `json:"_index,omitempty"`
	ID     string      `json:"_id,omitempty"`
	Score  float64     `json:"_score,omitempty"`
	Source interface{} `json:"_source,omitempty"`
}

type ElasticSearchHitsTotal struct {
	Value    float64 `json:"value,omitempty"`
	Relation string  `json:"relation,omitempty"`
}

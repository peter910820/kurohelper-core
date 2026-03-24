package vndb

// [VNDB]Request結構
//
// 對VNDB來講沒有必填項目，註解的必填項目是對於該專案來講的必填項目
// 所以預設值的部分可以完全不傳
//
// 這邊結構是根據需要的去對應，不是VNDB的完整結構
type BasicRequest struct {
	Filters           []any   `json:"filters"` // 必填
	Fields            string  `json:"fields"`  // 必填
	Sort              *string `json:"sort,omitempty"`
	Reverse           *bool   `json:"reverse,omitempty"`
	Results           *int    `json:"results,omitempty"`
	Page              *int    `json:"page,omitempty"`
	Count             *bool   `json:"count,omitempty"`
	CompactFilters    *bool   `json:"compact_filters,omitempty"`
	NormalizedFilters *bool   `json:"normalized_filters,omitempty"`
}

// vndb request factory
func VndbCreate() *BasicRequest {
	results := 100
	return &BasicRequest{
		Results: &results,
	}
}

func ptr(s string) *string { return &s }

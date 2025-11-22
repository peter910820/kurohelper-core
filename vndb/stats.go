package vndb

import (
	"encoding/json"
)

// 取得VNDB統計資料
func GetStats() (*Stats, error) {
	r, err := sendGetRequest("/stats")
	if err != nil {
		return nil, err
	}

	var res Stats
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

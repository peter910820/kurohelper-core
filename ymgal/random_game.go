package ymgal

import (
	"encoding/json"
)

type (
	randomGameResp struct {
		GID         int    `json:"gid"`
		DeveloperID int    `json:"developerId"`
		Name        string `json:"name"`
		ChineseName string `json:"chineseName"`
		HaveChinese bool   `json:"haveChinese"`
		MainImg     string `json:"mainImg"`
		ReleaseDate string `json:"releaseDate"`
		State       string `json:"state"`
	}
)

// 取得隨機遊戲
func GetRandomGame() ([]randomGameResp, error) {
	r, err := sendWithRetry(cfg.Endpoint + "/open/archive/random-game?num=1")
	if err != nil {
		return nil, err
	}

	var result basicResp[[]randomGameResp]

	err = json.Unmarshal(r, &result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, ErrAPIFailed{Code: result.Code}
	}

	return result.Data, nil
}

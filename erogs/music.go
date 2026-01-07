package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpercore "github.com/kuro-helper/kurohelper-core/v3"
)

func GetMusicByFuzzy(search string, idSearch bool) (*FuzzySearchMusicResponse, error) {
	searchJP := ""
	if !idSearch {
		searchJP = kurohelpercore.ZhTwToJp(search)
	}
	sql, err := buildFuzzySearchMusicSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchMusicResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

func GetMusicListByFuzzy(search string) (*[]FuzzySearchListResponse, error) {
	searchJP := kurohelpercore.ZhTwToJp(search)
	sql, err := buildFuzzySearchMusicListSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []FuzzySearchListResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

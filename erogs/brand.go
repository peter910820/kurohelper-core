package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpercore "github.com/kuro-helper/kurohelper-core/v3"
)

func GetBrandByFuzzy(search string) (*FuzzySearchBrandResponse, error) {
	searchJP := kurohelpercore.ZhTwToJp(search)
	sql, err := buildFuzzySearchBrandSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res FuzzySearchBrandResponse
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

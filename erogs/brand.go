package erogs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	Brand struct {
		ID            int    `json:"id"`
		BrandName     string `json:"brandname"`
		BrandFurigana string `json:"brandfurigana"`
		URL           string `json:"url"`
		Kind          string `json:"kind"`
		Lost          bool   `json:"lost"`
		DirectLink    bool   `json:"directlink"` // 網站可不可用
		Median        int    `json:"median"`     // 該品牌的遊戲評分中位數(一天更新一次)
		Twitter       string `json:"twitter"`
		Count2        int    `json:"count2"`
		CountAll      int    `json:"count_all"`
		Average2      int    `json:"average2"`
		Stdev         int    `json:"stdev"` // 標準偏差值(更新週期官方沒寫明確)
		GameList      []struct {
			ID       int    `json:"id"`
			GameName string `json:"gamename"`
			DMM      string `json:"dmm"` // dmm image
			Furigana string `json:"furigana"`
			SellDay  string `json:"sellday"`
			Model    string `json:"model"`
			Median   int    `json:"median"` // 分數中位數(一天更新一次)
			Stdev    int    `json:"stdev"`  // 分數標準偏差值(一天更新一次)
			Count2   int    `json:"count2"` // 分數計算的樣本數
			VNDB     string `json:"vndb"`   // *string
		} `json:"gamelist"`
	}
)

// Use erogs id search single brand data
func SearchBrandByID(id int) (*Brand, error) {
	sql := buildBrandSQL(fmt.Sprintf("WHERE id = '%d'", id))

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Brand
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

// Use kewords search single brand data
func SearchBrandByKeyword(keywords []string) (*Brand, error) {
	if keywords == nil {
		return nil, nil
	}

	// pre-build keySQL
	keySQL := "WHERE "
	var keywordSQLList []string
	for _, k := range keywords {
		formatK := buildSearchStringSQL(k)
		if strings.TrimSpace(formatK) != "" {
			keywordSQLList = append(keywordSQLList, fmt.Sprintf("brandname ILIKE '%s'", formatK))
		}
	}
	keySQL += strings.Join(keywordSQLList, " OR ")

	sql := buildBrandSQL(keySQL)

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Brand
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

// build search brand sql
// Arguments:
//   - keySQL: A pre-constructed SQL WHERE-clause fragment.
func buildBrandSQL(keySQL string) string {
	return fmt.Sprintf(`
WITH single_brand AS (
    SELECT
        id,
        brandname,
        brandfurigana,
        url,
        kind,
        lost,
        directlink,
        median,
        twitter,
        count2,
        count_all,
        average2,
        stdev
    FROM brandlist
    %s
    ORDER BY count2 DESC NULLS LAST, median DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(r)
FROM (
    SELECT 
        A.id, 
        A.brandname, 
        A.brandfurigana, 
        A.url, 
        A.kind, 
        A.lost, 
        A.directlink, 
        A.median, 
        A.twitter, 
        A.count2, 
        A.count_all, 
        A.average2, 
        A.stdev,
        (
            SELECT json_agg(
                json_build_object(
                    'id', g.id,
                    'gamename', g.gamename,
					'dmm', g.dmm,
                    'furigana', g.furigana,
                    'sellday', g.sellday,
                    'median', g.median,
                    'model', g.model,
                    'stdev', g.stdev,
                    'count2', g.count2,
                    'vndb', g.vndb
                ) ORDER BY g.sellday DESC
            )
            FROM gamelist g
            WHERE g.brandname = A.id
        ) AS gamelist
    FROM single_brand A
) r;
`, keySQL)
}

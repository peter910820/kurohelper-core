package erogs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 只抓一筆(LIMIT 1)
type Creator struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	TwitterUsername string `json:"twitter_username"`
	Blog            string `json:"blog"`
	Pixiv           *int   `json:"pixiv"`
	Games           []struct {
		Gamename string `json:"gamename"`
		SellDay  string `json:"sellday"`
		Median   int    `json:"median"`
		CountAll int    `json:"count2"`
		Shokushu []struct {
			Shubetu           int    `json:"shubetu"`
			ShubetuDetail     int    `json:"shubetu_detail"`
			ShubetuDetailName string `json:"shubetu_detail_name"` // *string
		} `json:"shokushu"` // 有可能一個遊戲有多種身分
	} `json:"games"` // 參與過的遊戲
}

type CreatorList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Use kewords search creator list data
func SearchCreatorListByKeyword(keywords []string) ([]CreatorList, error) {
	if keywords == nil {
		return nil, nil
	}

	// pre-build keySQL
	keySQL := "WHERE "
	var keywordSQLList []string
	for _, k := range keywords {
		formatK := buildSearchStringSQL(k)
		if strings.TrimSpace(formatK) != "" {
			keywordSQLList = append(keywordSQLList, fmt.Sprintf("cr.name ILIKE '%s'", formatK))
		}
	}
	keySQL += strings.Join(keywordSQLList, " OR")

	sql := buildCreatorListSQL(keySQL)

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []CreatorList
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return res, nil
}

// Use erogs id search single creator data
func SearchCreatorByID(id int) (*Creator, error) {
	sql := buildCreatorSQL(fmt.Sprintf("WHERE cr.id = '%d'", id))

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Creator
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Use kewords search single creator data
func SearchCreatorByKeyword(keywords []string) (*Creator, error) {
	if keywords == nil {
		return nil, nil
	}

	// pre-build keySQL
	keySQL := "WHERE "
	var keywordSQLList []string
	for _, k := range keywords {
		formatK := buildSearchStringSQL(k)
		if strings.TrimSpace(formatK) != "" {
			keywordSQLList = append(keywordSQLList, fmt.Sprintf("cr.name ILIKE '%s'", formatK))
		}
	}
	keySQL += strings.Join(keywordSQLList, " OR")

	sql := buildCreatorSQL(keySQL)

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Creator
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// build search creator list sql
// Arguments:
//   - keySQL: A pre-constructed SQL WHERE-clause fragment.
func buildCreatorListSQL(keySQL string) string {
	return fmt.Sprintf(`
SELECT json_agg(row_to_json(c))
FROM (
    SELECT
        cr.id,
        cr.name
    FROM createrlist cr
    %s
    LIMIT 200
) AS c;
`, keySQL)
}

// build search creator sql
// Arguments:
//   - keySQL: A pre-constructed SQL WHERE-clause fragment.
func buildCreatorSQL(keySQL string) string {
	return fmt.Sprintf(`
SELECT row_to_json(c)
FROM (
    SELECT
        cr.id,
        cr.name,
        cr.twitter_username,
        cr.blog,
        cr.pixiv,
        (
            SELECT json_agg(game_data)
            FROM (
                SELECT
                    g.gamename,
                    g.sellday,
                    g.median,
                    g.count2,
                    (
                        SELECT json_agg(
                            json_build_object(
                                'shubetu', s2.shubetu,
                                'shubetu_detail', s2.shubetu_detail,
                                'shubetu_detail_name', s2.shubetu_detail_name
                            )
                        )
                        FROM shokushu s2
                        WHERE s2.creater = cr.id
                          AND s2.game = g.id
                    ) AS shokushu
                FROM gamelist g
                WHERE EXISTS (
                    SELECT 1 
                    FROM shokushu s3
                    WHERE s3.creater = cr.id
                      AND s3.game = g.id
                )
                GROUP BY g.id, g.gamename, g.sellday, g.median, g.count2
            ) AS game_data
        ) AS games
    FROM createrlist cr
    %s
    LIMIT 1
) AS c;`, keySQL)
}

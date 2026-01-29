package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpercore "kurohelper-core"
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

func GetCreatorByFuzzy(search string, idSearch bool) (*Creator, error) {
	searchJP := ""
	if !idSearch {
		searchJP = kurohelpercore.ZhTwToJp(search)
	}
	sql, err := buildFuzzySearchCreatorSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

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

func GetCreatorListByFuzzy(search string) (*[]CreatorList, error) {
	searchJP := kurohelpercore.ZhTwToJp(search)
	sql, err := buildFuzzySearchCreatorListSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

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

	return &res, nil
}

func buildFuzzySearchCreatorSQL(searchTW string, searchJP string, idSearch bool) (string, error) {
	searchString := ""
	if idSearch {
		idString := searchTW[1:]
		searchString = fmt.Sprintf("WHERE cr.id = %s", idString)
	} else {
		resultTW, err := buildSearchStringSQL(searchTW)
		if err != nil {
			return "", err
		}

		resultJP, err := buildSearchStringSQL(searchJP)
		if err != nil {
			return "", err
		}
		searchString = fmt.Sprintf("WHERE cr.name ILIKE '%s' OR cr.name ILIKE '%s'", resultTW, resultJP)
	}
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
) AS c;`, searchString), nil
}

func buildFuzzySearchCreatorListSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
SELECT json_agg(row_to_json(c))
FROM (
    SELECT
        cr.id,
        cr.name
    FROM createrlist cr
    WHERE cr.name ILIKE '%s' OR cr.name ILIKE '%s'
    LIMIT 200
) AS c;
`, resultTW, resultJP), nil
}

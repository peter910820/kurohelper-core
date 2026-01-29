package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpercore "kurohelper-core"
)

type Game struct {
	ID                               int    `json:"id"`
	BrandID                          int    `json:"brandid"`
	BrandName                        string `json:"brandname"`
	Gamename                         string `json:"gamename"`
	SellDay                          string `json:"sellday"`
	Model                            string `json:"model"`
	DMM                              string `json:"dmm"` // dmm image
	Median                           string `json:"median"`
	TokutenCount                     string `json:"count2"`
	TotalPlayTimeMedian              string `json:"total_play_time_median"`
	TimeBeforeUnderstandingFunMedian string `json:"time_before_understanding_fun_median"`
	Okazu                            string `json:"okazu"`
	Erogame                          string `json:"erogame"`
	Genre                            string `json:"genre"`
	BannerUrl                        string `json:"banner_url"`
	SteamId                          string `json:"steam"`
	VndbId                           string `json:"vndb"`
	Shoukai                          string `json:"shoukai"`
	Junni                            int    `json:"junni"`
	CreatorShubetu                   []struct {
		ShubetuType       int    `json:"shubetu_type"`
		CreatorName       string `json:"creater_name"`
		ShubetuDetailType int    `json:"shubetu_detail_type"`
		ShubetuDetailName string `json:"shubetu_detail_name"`
	} `json:"shubetu_detail"`
}

type GameList struct {
	ID                               int    `json:"id"`
	Name                             string `json:"name"`
	Category                         string `json:"category"`
	DMM                              string `json:"dmm"` // dmm image
	Median                           string `json:"median"`
	TokutenCount                     string `json:"count2"`
	TotalPlayTimeMedian              string `json:"total_play_time_median"`
	TimeBeforeUnderstandingFunMedian string `json:"time_before_understanding_fun_median"`
}

func GetGameByFuzzy(search string, idSearch bool) (*Game, error) {
	searchJP := ""
	if !idSearch {
		searchJP = kurohelpercore.ZhTwToJp(search)
	}
	sql, err := buildFuzzySearchGameSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Game
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func GetGameListByFuzzy(search string) (*[]GameList, error) {
	searchJP := kurohelpercore.ZhTwToJp(search)
	sql, err := buildFuzzySearchGameListSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []GameList
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

func buildFuzzySearchGameSQL(searchTW string, searchJP string, idSearch bool) (string, error) {
	searchString := ""
	if idSearch {
		idString := searchTW[1:]
		searchString = fmt.Sprintf("WHERE id = %s", idString)
	} else {
		resultTW, err := buildSearchStringSQL(searchTW)
		if err != nil {
			return "", err
		}

		resultJP, err := buildSearchStringSQL(searchJP)
		if err != nil {
			return "", err
		}
		searchString = fmt.Sprintf("WHERE gamename ILIKE '%s' OR gamename ILIKE '%s'", resultTW, resultJP)
	}
	return fmt.Sprintf(`
WITH filtered_games AS (
    SELECT *
    FROM gamelist
    %s
    ORDER BY count2 DESC NULLS LAST, median DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(t)
FROM (
    SELECT g.id,
           b.id AS brandid, 
           b.brandname,
           g.gamename,
           g.sellday,
           g.model,
           g.dmm,
           COALESCE(g.median::text, '無') AS median,
           COALESCE(g.count2::text, '無') AS count2,
           COALESCE(g.total_play_time_median::text, '無') AS total_play_time_median,
           COALESCE(g.time_before_understanding_fun_median::text, '無') AS time_before_understanding_fun_median,
           COALESCE(g.okazu::text, '未收錄') AS okazu,
           COALESCE(g.erogame::text, '未收錄') AS erogame,
           COALESCE(g.banner_url, '') AS banner_url,
           COALESCE(g.genre, '無') AS genre,
           COALESCE(g.steam::text, '') AS steam,
           COALESCE(g.vndb, '') AS vndb,
           j.junni,
           g.shoukai,
           s.shubetu_detail
    FROM filtered_games g
    LEFT JOIN LATERAL (
        SELECT json_agg(
                   json_build_object(
                       'shubetu_type', s.shubetu,
                       'creater_name', c.name,
                       'shubetu_detail_type', s.shubetu_detail,
                       'shubetu_detail_name', s.shubetu_detail_name
                   )
               ) AS shubetu_detail
        FROM shokushu s
        LEFT JOIN createrlist c ON c.id = s.creater
        WHERE s.game = g.id AND s.shubetu != 7
    ) s ON TRUE
    LEFT JOIN brandlist b ON b.id = g.brandname
    LEFT JOIN LATERAL (
        SELECT j.junni
        FROM junirirekimedian j
        WHERE j.game = g.id
        ORDER BY j.date DESC NULLS LAST
        LIMIT 1
    ) j ON TRUE
) t;
`, searchString), nil
}

func buildFuzzySearchGameListSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
SELECT json_agg(row_to_json(t))
FROM (
    SELECT g.id,
           g.gamename AS name,
           g.model AS category,
           g.dmm,
           COALESCE(g.median::text, '無') AS median,
           COALESCE(g.count2::text, '無') AS count2,
           COALESCE(g.total_play_time_median::text, '無') AS total_play_time_median,
           COALESCE(g.time_before_understanding_fun_median::text, '無') AS time_before_understanding_fun_median
    FROM gamelist g
    WHERE gamename ILIKE '%s' OR gamename ILIKE '%s'
    ORDER BY count2 DESC NULLS LAST, median DESC NULLS LAST
    LIMIT 200
) t;
`, resultTW, resultJP), nil
}

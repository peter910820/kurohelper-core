package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpercore "github.com/kuro-helper/kurohelper-core/v3"
)

type (
	MusicList struct {
		ID           int    `json:"id"`            // 歌曲ID
		Name         string `json:"name"`          // 歌曲名稱
		TokutenCount string `json:"tokuten_count"` // 樣本數
		AvgTokuten   string `json:"avg_tokuten"`   // 平均分數
		Category     string `json:"category"`      // 歌曲類型
		GameName     string `json:"game_name"`     // 對應的遊戲名稱
		GameDMM      string `json:"game_dmm"`      // 對應的DMM圖片(代號)
	}

	Music struct {
		ID             int                 `json:"music_id"`         // 歌曲ID
		MusicName      string              `json:"musicname"`        // 歌曲名稱
		PlayTime       string              `json:"playtime"`         // 歌曲長度
		ReleaseDate    string              `json:"releasedate"`      // 發售日
		AvgTokuten     float64             `json:"avg_tokuten"`      // 平均分數
		TokutenCount   int                 `json:"tokuten_count"`    // 樣本數
		Singers        string              `json:"singer_name"`      // 歌手
		Lyrics         string              `json:"lyric_name"`       // 作詞家
		Arrangments    string              `json:"arrangement_name"` // 作曲家
		Compositions   string              `json:"composition_name"` // 編曲家
		GameCategories []MusicGameCategory `json:"game_categories"`  // 遊戲資料相關
		Album          string              `json:"album_name"`       // 收錄的專輯
	}

	// 查詢音樂的對應遊戲資料
	MusicGameCategory struct {
		GameDMM   string `json:"dmm"`        // 對應的DMM圖片(代號)
		Category  string `json:"category"`   // 歌曲類型
		GameName  string `json:"game_name"`  // 遊戲名稱
		GameModel string `json:"game_model"` // 遊戲平台
	}
)

func GetMusicListByFuzzy(search string) ([]MusicList, error) {
	searchJP := kurohelpercore.ZhTwToJp(search)
	sql, err := buildMusicListSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []MusicList
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return res, nil
}

func GetMusicByFuzzy(search string, idSearch bool) (*Music, error) {
	searchJP := ""
	if !idSearch {
		searchJP = kurohelpercore.ZhTwToJp(search)
	}
	sql, err := buildMusicSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Music
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

func buildMusicListSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
WITH filtered_music AS (
    SELECT 
        m.id AS music_id,
        m.name AS musicname,
        ROUND(AVG(LEAST(ut.tokuten, 100))::numeric, 2) AS avg_tokuten,
        COUNT(DISTINCT ut.uid) AS tokuten_count
    FROM musiclist m
    LEFT JOIN usermusic_tokuten ut ON ut.music = m.id
    WHERE m.name ILIKE '%s' OR m.name ILIKE '%s'
    GROUP BY m.id, m.name, m.playtime, m.releasedate
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
    LIMIT 200
)
SELECT json_agg(row_to_json(t))
FROM (
    SELECT 
        m.music_id AS id,
        m.musicname AS name,
        m.tokuten_count,
        m.avg_tokuten,
        STRING_AGG(DISTINCT gm.category, ',') AS category,
        STRING_AGG(DISTINCT gmlist.gamename, ', ') AS game_name,
        STRING_AGG(DISTINCT gmlist.dmm, ', ') AS game_dmm
    FROM filtered_music m
    LEFT JOIN game_music gm ON gm.music = m.music_id,
    LEFT JOIN gamelist gmlist ON gmlist.id = gm.game,
    GROUP BY m.music_id,m.musicname,m.tokuten_count,m.avg_tokuten
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
) t;
`, resultTW, resultJP), nil
}

func buildMusicSQL(searchTW string, searchJP string, idSearch bool) (string, error) {
	searchString := ""
	if idSearch {
		idString := searchTW[1:]
		searchString = fmt.Sprintf("WHERE m.id = %s", idString)
	} else {
		resultTW, err := buildSearchStringSQL(searchTW)
		if err != nil {
			return "", err
		}

		resultJP, err := buildSearchStringSQL(searchJP)
		if err != nil {
			return "", err
		}
		searchString = fmt.Sprintf("WHERE m.name ILIKE '%s' OR m.name ILIKE '%s'", resultTW, resultJP)
	}
	return fmt.Sprintf(`
WITH filtered_music AS (
    SELECT 
        m.id AS music_id,
        m.name AS musicname,
        m.playtime,
        m.releasedate,
        ROUND(AVG(LEAST(ut.tokuten, 100))::numeric, 2) AS avg_tokuten,
        COUNT(DISTINCT ut.uid) AS tokuten_count
    FROM musiclist m
    LEFT JOIN usermusic_tokuten ut ON ut.music = m.id
    %s
    GROUP BY m.id, m.name, m.playtime, m.releasedate
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(t)
FROM (
    SELECT 
        m.music_id,
        m.musicname,
        m.playtime,
        m.releasedate,
        m.avg_tokuten,
        m.tokuten_count,
        COALESCE(STRING_AGG(DISTINCT s_c.name, ','), '無') AS singer_name,
        COALESCE(STRING_AGG(DISTINCT l_c.name, ','), '無') AS lyric_name,
        COALESCE(STRING_AGG(DISTINCT a_c.name, ','), '無') AS arrangement_name,
        COALESCE(STRING_AGG(DISTINCT comp_c.name, ','), '無') AS composition_name,
        json_agg(
            DISTINCT jsonb_build_object(
                'game_name', g.gamename,
                'game_model', g.model,
                'dmm', g.dmm,
                'category', gm.category
            )
        ) AS game_categories,
        COALESCE(STRING_AGG(DISTINCT mi.name, ','), '') AS album_name
    FROM filtered_music m
    LEFT JOIN singer s ON s.music = m.music_id
    LEFT JOIN createrlist s_c ON s_c.id = s.creater
    LEFT JOIN lyrics l ON l.music = m.music_id
    LEFT JOIN createrlist l_c ON l_c.id = l.creater
    LEFT JOIN arrangement a ON a.music = m.music_id
    LEFT JOIN createrlist a_c ON a_c.id = a.creater
    LEFT JOIN composition comp ON comp.music = m.music_id
    LEFT JOIN createrlist comp_c ON comp_c.id = comp.creater
    LEFT JOIN game_music gm ON gm.music = m.music_id
    LEFT JOIN gamelist g ON g.id = gm.game
    LEFT JOIN musicitem_music mim ON mim.music = m.music_id
    LEFT JOIN musicitemlist mi ON mi.id = mim.musicitem
    GROUP BY m.music_id, m.musicname, m.playtime, m.releasedate, m.avg_tokuten, m.tokuten_count
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
) t;
`, searchString), nil
}

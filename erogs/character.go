package erogs

import (
	"encoding/json"
	"fmt"

	kurohelpercore "kurohelper-core"
)

type Character struct {
	ID            int    `json:"id"`
	CharacterName string `json:"name"`
	Sex           string `json:"sex"`
	BloodType     string `json:"bloodtype"`
	Birthday      string `json:"birthday"`
	GameName      string `json:"gamename"`
	URL           string `json:"url"`
	FormalExplain string `json:"formal_explanation"`
	Age           string `json:"age"`
	Bust          string `json:"bust"`
	Waist         string `json:"waist"`
	Hip           string `json:"hip"`
	Height        string `json:"height"`
	Weight        string `json:"weight"`
	Cup           string `json:"cup"`
	Role          int    `json:"role"`
	CreatorName   string `json:"creater_name"`
}

type CharacterList struct {
	ID       int    `json:"id"`
	GameName string `json:"gamename"`
	Category string `json:"category"`
	Model    string `json:"model"`
}

func GetCharacterByFuzzy(search string, idSearch bool) (*Character, error) {
	searchJP := ""
	if !idSearch {
		searchJP = kurohelpercore.ZhTwToJp(search)
	}
	sql, err := buildFuzzySearchCharacterSQL(search, searchJP, idSearch)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res Character
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func GetCharacterListByFuzzy(search string) (*[]CharacterList, error) {
	searchJP := kurohelpercore.ZhTwToJp(search)
	sql, err := buildFuzzySearchCharacterListSQL(search, searchJP)
	if err != nil {
		return nil, err
	}

	jsonText, err := sendPostRequest(sql)
	if err != nil {
		return nil, err
	}

	var res []CharacterList
	err = json.Unmarshal([]byte(jsonText), &res)
	if err != nil {
		fmt.Println(jsonText)
		return nil, err
	}

	return &res, nil
}

func buildFuzzySearchCharacterSQL(searchTW string, searchJP string, idSearch bool) (string, error) {
	searchString := ""
	if idSearch {
		idString := searchTW[1:]
		searchString = fmt.Sprintf("WHERE ch.id = %s", idString)
	} else {
		resultTW, err := buildSearchStringSQL(searchTW)
		if err != nil {
			return "", err
		}

		resultJP, err := buildSearchStringSQL(searchJP)
		if err != nil {
			return "", err
		}
		searchString = fmt.Sprintf("WHERE ch.name ILIKE '%s' OR ch.name ILIKE '%s'", resultTW, resultJP)
	}
	return fmt.Sprintf(`
WITH filtered_character AS (
    SELECT 
        ch.id AS char_id,
        ch.name,
        ch.sex,
        ch.bloodtype,
        ch.birthday,
        g.gamename,
        g.count2,
        g.median,
        a.role,
        a.url,
        a.formal_explanation,
        a.age,
        a.bust,
        a.waist,
        a.hip,
        a.height,
        a.weight,
        a.cup
    FROM characterlist ch
    LEFT JOIN appearance a ON a.character = ch.id
    LEFT JOIN gamelist g ON g.id = a.game
    %s
    ORDER BY g.count2 DESC NULLS LAST, g.median DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(t)
FROM (
    SELECT ch.char_id,
        ch.name,
        ch.sex,
        ch.bloodtype,
        ch.birthday,
        ch.gamename,
        ch.url,
        COALESCE(NULLIF(ch.formal_explanation::text, ''), '無') AS formal_explanation,
        COALESCE(NULLIF(ch.age::text, ''), '未收錄') AS age,
        ch.bust,
        ch.waist,
        ch.hip,
        COALESCE(NULLIF(ch.height::text, ''), '未收錄') AS height,
        COALESCE(NULLIF(ch.weight::text, ''), '未收錄') AS weight,
        COALESCE(NULLIF(ch.cup::text, ''), '未收錄') AS cup,
        chrole.id AS role,
        cr.name AS creater_name
    FROM filtered_character ch
    LEFT JOIN character_rolelist chrole ON chrole.id = ch.role
    LEFT JOIN appearance_actor a ON a.character = ch.char_id
    LEFT JOIN createrlist cr ON cr.id = a.actor
) t;
`, searchString), nil
}

func buildFuzzySearchCharacterListSQL(searchTW string, searchJP string) (string, error) {
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
        SELECT 
            ch.id AS id,
            ch.name,
            g.gamename AS gamename,
            g.model
        FROM characterlist ch
        LEFT JOIN appearance a ON a.character = ch.id
        LEFT JOIN gamelist g ON g.id = a.game
        WHERE ch.name ILIKE '%s' OR ch.name ILIKE '%s'
        ORDER BY g.count2 DESC NULLS LAST, g.median DESC NULLS LAST
        LIMIT 200
    ) t;
    `, resultTW, resultJP), nil
}

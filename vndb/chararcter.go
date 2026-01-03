package vndb

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	kurohelpercore "github.com/kuro-helper/core/v2"
)

// 取得VNDB角色(模糊搜尋)
func GetCharacterByFuzzy(keyword string) (*CharacterSearchResponse, error) {
	reqCharacter := VndbCreate() // 建立基本request結構

	// 依照關鍵字的相關度排序
	reqCharacterSort := "searchrank"
	reqCharacter.Sort = &reqCharacterSort

	// 限制回傳一筆結果
	reqCharacterResults := 1
	reqCharacter.Results = &reqCharacterResults

	// 指定要取得的欄位
	basicFields := "id, name, original, aliases, description, image.url, blood_type, height, weight, bust, waist, hips, cup, age, birthday, sex, gender"
	vnsFields := "vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main"
	allFields := []string{
		basicFields,
		vnsFields,
	}
	reqCharacter.Fields = strings.Join(allFields, ", ")

	// 設定搜尋條件
	reqCharacter.Filters = []any{"search", "=", keyword}

	jsonCharacter, err := json.Marshal(reqCharacter)
	if err != nil {
		return nil, err
	}

	r, err := sendPostRequest("/character", jsonCharacter)
	if err != nil {
		return nil, err
	}

	var resCharacters BasicResponse[CharacterSearchResponse]
	err = json.Unmarshal(r, &resCharacters)
	if err != nil {
		return nil, err
	}
	if len(resCharacters.Results) == 0 {
		return nil, kurohelpercore.ErrSearchNoContent
	}

	// 取得角色詳細資料
	err = GetCharacterDetail(resCharacters.Results[0].ID, &resCharacters)
	if err != nil {
		return nil, err
	}
	return &resCharacters.Results[0], nil
}

// 用VNDB角色ID取得VNDB角色
func GetCharacterByID(keyword string) (*CharacterSearchResponse, error) {
	reqCharacter := VndbCreate() // 建立基本request結構

	// 不需要排序
	reqCharacterSort := ""
	reqCharacter.Sort = &reqCharacterSort

	// 限制回傳一筆結果
	reqCharacterResults := 1
	reqCharacter.Results = &reqCharacterResults

	// 指定要取得的欄位
	basicFields := "id, name, original, aliases, description, image.url, blood_type, height, weight, bust, waist, hips, cup, age, birthday, sex, gender"
	vnsFields := "vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main"
	allFields := []string{
		basicFields,
		vnsFields,
	}
	reqCharacter.Fields = strings.Join(allFields, ", ")

	// 設定搜尋條件
	reqCharacter.Filters = []any{"id", "=", keyword}

	jsonCharacter, err := json.Marshal(reqCharacter)
	if err != nil {
		return nil, err
	}

	r, err := sendPostRequest("/character", jsonCharacter)
	if err != nil {
		return nil, err
	}

	var resCharacters BasicResponse[CharacterSearchResponse]
	err = json.Unmarshal(r, &resCharacters)
	if err != nil {
		return nil, err
	}
	if len(resCharacters.Results) == 0 {
		return nil, kurohelpercore.ErrSearchNoContent
	}

	// 取得角色詳細資料
	err = GetCharacterDetail(resCharacters.Results[0].ID, &resCharacters)
	if err != nil {
		return nil, err
	}
	return &resCharacters.Results[0], nil
}

// 取得VNDB隨機角色
func GetRandomCharacter(opt string) (*CharacterSearchResponse, error) {
	reqCharacter := VndbCreate() // 建立基本request結構

	// 不需要排序
	reqCharacterSort := ""
	reqCharacter.Sort = &reqCharacterSort

	// 限制回傳一筆結果
	reqCharacterResults := 1
	reqCharacter.Results = &reqCharacterResults

	// 根據角色身分過濾結果
	reqCharacter.Filters = []any{"and", []any{"vn", "=", []any{"votecount", ">=", "50"}}}
	switch opt {
	case "":
		fallthrough
	case "1":
		reqCharacter.Filters = append(reqCharacter.Filters, []any{"or", []any{"role", "=", "main"}, []any{"role", "=", "primary"}}) // 主角
	case "2":
		reqCharacter.Filters = append(reqCharacter.Filters, []any{"or", []any{"role", "=", "side"}, []any{"role", "=", "appear"}}) // 配角
	}

	// 指定要取得的欄位
	basicFields := "id, name, original, aliases, description, image.url, blood_type, height, weight, bust, waist, hips, cup, age, birthday, sex, gender"
	vnsFields := "vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main"
	allFields := []string{
		basicFields,
		vnsFields,
	}
	reqCharacter.Fields = strings.Join(allFields, ", ")

	// 設定搜尋條件
	resStat, err := GetStats() // 獲取角色id總數
	if err != nil {
		return nil, err
	}
	var resCharacters BasicResponse[CharacterSearchResponse]
	var randomCharacterID string
	for {
		randomCharacterID = fmt.Sprintf("c%d", rand.Intn(resStat.Chars))
		reqCharacter.Filters = append(reqCharacter.Filters, []any{"and", []any{"id", ">=", randomCharacterID}, []any{"vn", "=", []any{"votecount", ">=", "100"}}})

		jsonCharacter, err := json.Marshal(reqCharacter)
		if err != nil {
			return nil, err
		}

		r, err := sendPostRequest("/character", jsonCharacter)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(r, &resCharacters)
		if err != nil {
			return nil, err
		}

		if len(resCharacters.Results) != 0 {
			break
		}
	}
	// 取得角色詳細資料
	err = GetCharacterDetail(resCharacters.Results[0].ID, &resCharacters)
	if err != nil {
		return nil, err
	}

	return &resCharacters.Results[0], nil
}

func GetCharacterDetail(characterID string, resCharacters *BasicResponse[CharacterSearchResponse]) error {
	reqVn := VndbCreate()

	characterIDFilter := []any{"id", "=", characterID}
	reqVn.Filters = []any{"character", "=", characterIDFilter}
	reqVn.Fields = "va.staff.name, va.staff.original, va.character.id"

	jsonVn, err := json.Marshal(reqVn)
	if err != nil {
		return err
	}

	r, err := sendPostRequest("/vn", jsonVn)
	if err != nil {
		return err
	}

	var resVn BasicResponse[GetVnUseIDResponse]
	err = json.Unmarshal(r, &resVn)
	if err != nil {
		return err
	}

	var vasMap = make(map[string]bool) // 去重
	var vas []string
	for _, vn := range resVn.Results {
		for _, va := range vn.Va {
			if va.Character.ID == characterID {
				if va.Staff.Original != "" {
					vasMap[va.Staff.Original] = true
				} else {
					vasMap[va.Staff.Name] = true
				}
			}
		}
	}

	if len(vasMap) == 0 {
		resCharacters.Results[0].Vas = []string{"未收錄"}
	} else {
		for va := range vasMap {
			vas = append(vas, va)
		}
		resCharacters.Results[0].Vas = vas
	}

	return nil
}

// 取得VNDB角色列表(模糊搜尋)
func GetCharacterListByFuzzy(keyword string) (*[]CharacterSearchResponse, error) {
	reqCharacter := VndbCreate()
	reqCharacter.Filters = []any{"search", "=", keyword}
	reqCharacterSort := "searchrank"
	reqCharacter.Sort = &reqCharacterSort
	basicFields := "id, name, original"
	vnsFields := "vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main"
	allFields := []string{
		basicFields,
		vnsFields,
	}
	reqCharacter.Fields = strings.Join(allFields, ", ")
	jsonCharacter, err := json.Marshal(reqCharacter)
	if err != nil {
		return nil, err
	}
	r, err := sendPostRequest("/character", jsonCharacter)
	if err != nil {
		return nil, err
	}
	var resCharacters BasicResponse[CharacterSearchResponse]
	err = json.Unmarshal(r, &resCharacters)
	if err != nil {
		return nil, err
	}

	return &resCharacters.Results, nil
}

func ConvertBBCodeToMarkdown(text string) string {
	// 1. 處理配對的 URL 標籤
	reURL := regexp.MustCompile(`\[url=(.+?)\](.+?)\[/url\]`)
	text = reURL.ReplaceAllString(text, "[$2]($1)")

	// 2. 處理配對的 spoiler 標籤（支援多行）
	reSpoiler := regexp.MustCompile(`(?s)\[spoiler\](.+?)\[/spoiler\]`)
	text = reSpoiler.ReplaceAllString(text, "||$1||")

	// 3. 清理未配對的殘留標籤
	text = strings.ReplaceAll(text, "[spoiler]", "")
	text = strings.ReplaceAll(text, "[/spoiler]", "")

	// 4. 將角色ID轉換成連結[Sara](/c40662)
	reCharacterID := regexp.MustCompile(`\[(.+?)\]\(/c(\d+?)\)`)
	text = reCharacterID.ReplaceAllString(text, "[$1](https://vndb.org/c$2)")
	return strings.TrimSpace(text)
}

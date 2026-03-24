package vndb

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	kurohelpercore "kurohelper-core"
)

var characterFields = []string{
	// basic fields
	"id, name, original, aliases, description, image.url, blood_type, height, weight, bust, waist, hips, cup, age, birthday, sex, gender",
	// vns fields
	"vns.title, vns.alttitle, vns.spoiler, vns.role, vns.titles.title, vns.titles.main",
}

var (
	reURL         = regexp.MustCompile(`\[url=(.+?)\](.+?)\[/url\]`)
	reSpoiler     = regexp.MustCompile(`(?s)\[spoiler\](.+?)\[/spoiler\]`)
	reCharacterID = regexp.MustCompile(`\[(.+?)\]\(/c(\d+?)\)`)
)

// 取得 VNDB 角色（模糊搜尋）
func GetCharacterByFuzzy(keyword string) (*CharacterSearchResponse, error) {
	sort := "searchrank"
	return fetchCharacter([]any{"search", "=", keyword}, sort)
}

// 用 VNDB 角色 ID 取得角色
func GetCharacterByID(id string) (*CharacterSearchResponse, error) {
	return fetchCharacter([]any{"id", "=", id}, "")
}

// 取得 VNDB 隨機角色
func GetRandomCharacter(opt string) (*CharacterSearchResponse, error) {
	resStat, err := GetStats()
	if err != nil {
		return nil, err
	}

	baseFilter := []any{"vn", "=", []any{"and", []any{"votecount", ">=", "30"}, []any{"rating", ">=", "70"}}}
	var roleFilter []any
	switch opt {
	case "", "1":
		roleFilter = []any{"or", []any{"role", "=", "main"}, []any{"role", "=", "primary"}}
	case "2":
		roleFilter = []any{"or", []any{"role", "=", "side"}, []any{"role", "=", "appear"}}
	default:
		roleFilter = []any{"or", []any{"role", "=", "main"}, []any{"role", "=", "primary"}}
	}

	for range 3 {
		randomID := fmt.Sprintf("c%d", rand.Intn(resStat.Chars))
		idFilter := []any{"and", []any{"id", ">=", randomID}, []any{"vn", "=", []any{"votecount", ">=", "100"}}}
		filters := []any{"and", baseFilter, roleFilter, idFilter}

		req := VndbCreate()
		req.Sort = ptr("")
		reqResults := 1
		req.Results = &reqResults
		req.Fields = strings.Join(characterFields, ", ")
		req.Filters = filters

		body, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		r, err := sendPostRequest("/character", body)
		if err != nil {
			return nil, err
		}

		var res BasicResponse[CharacterSearchResponse]
		if err = json.Unmarshal(r, &res); err != nil {
			return nil, err
		}
		if len(res.Results) == 0 {
			continue
		}

		if err = getCharacterDetail(res.Results[0].ID, &res); err != nil {
			return nil, err
		}
		return &res.Results[0], nil
	}
	return nil, kurohelpercore.ErrSearchNoContent
}

// 取得 VNDB 角色列表（模糊搜尋）
func GetCharacterListByFuzzy(keyword string) ([]CharacterSearchResponse, error) {
	req := VndbCreate()
	req.Filters = []any{"search", "=", keyword}
	sort := "searchrank"
	req.Sort = &sort
	req.Fields = strings.Join(characterFields, ", ")

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := sendPostRequest("/character", body)
	if err != nil {
		return nil, err
	}

	var res BasicResponse[CharacterSearchResponse]
	if err = json.Unmarshal(r, &res); err != nil {
		return nil, err
	}
	return res.Results, nil
}

// 內部共用：依 filters 查詢單一角色並取得詳細資料
func fetchCharacter(filters []any, sort string) (*CharacterSearchResponse, error) {
	req := VndbCreate()
	req.Sort = &sort
	reqResults := 1
	req.Results = &reqResults
	req.Fields = strings.Join(characterFields, ", ")
	req.Filters = filters

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := sendPostRequest("/character", body)
	if err != nil {
		return nil, err
	}

	var res BasicResponse[CharacterSearchResponse]
	if err = json.Unmarshal(r, &res); err != nil {
		return nil, err
	}
	if len(res.Results) == 0 {
		return nil, kurohelpercore.ErrSearchNoContent
	}

	if err = getCharacterDetail(res.Results[0].ID, &res); err != nil {
		return nil, err
	}
	return &res.Results[0], nil
}

// 取得角色 VA 資訊並寫入 resCharacters(會額外呼叫 /vn API)
func getCharacterDetail(characterID string, resCharacters *BasicResponse[CharacterSearchResponse]) error {
	req := VndbCreate()
	req.Filters = []any{"character", "=", []any{"id", "=", characterID}}
	req.Fields = "va.staff.name, va.staff.original, va.character.id"

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	r, err := sendPostRequest("/vn", body)
	if err != nil {
		return err
	}

	var resVn BasicResponse[GetVnUseIDResponse]
	if err = json.Unmarshal(r, &resVn); err != nil {
		return err
	}

	vasMap := make(map[string]struct{})
	for _, vn := range resVn.Results {
		for _, va := range vn.Va {
			if va.Character.ID != characterID {
				continue
			}
			if va.Staff.Original != "" {
				vasMap[va.Staff.Original] = struct{}{}
			} else {
				vasMap[va.Staff.Name] = struct{}{}
			}
		}
	}

	if len(vasMap) == 0 {
		resCharacters.Results[0].Vas = []string{"未收錄"}
	} else {
		vas := make([]string, 0, len(vasMap))
		for va := range vasMap {
			vas = append(vas, va)
		}
		resCharacters.Results[0].Vas = vas
	}
	return nil
}

// 將 BBCode 轉為 Markdown
func ConvertBBCodeToMarkdown(text string) string {
	text = reURL.ReplaceAllString(text, "[$2]($1)")
	text = reSpoiler.ReplaceAllString(text, "||$1||")
	text = strings.ReplaceAll(text, "[spoiler]", "")
	text = strings.ReplaceAll(text, "[/spoiler]", "")
	text = reCharacterID.ReplaceAllString(text, "[$1](https://vndb.org/c$2)")
	return strings.TrimSpace(text)
}

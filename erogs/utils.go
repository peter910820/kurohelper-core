package erogs

import (
	"fmt"
	kurohelpercore "kurohelper-core"
	"strings"
)

var (
	ShubetuMap = map[int]string{
		1: "原画",
		2: "シナリオ",
		3: "音楽",
		4: "キャラデザ",
		5: "声優",
		6: "歌手",
		7: "その他",
	}
	Role = map[int]string{
		1: "メイン",
		2: "サブ",
		3: "主人公",
		4: "その他",
	}
)

func MakeDMMImageURL(dmm string) string {
	return fmt.Sprintf("https://pics.dmm.co.jp/digital/pcgame/%[1]s/%[1]spl.jpg", dmm)
}

func buildSearchStringSQL(search string) (string, error) {
	search = strings.ReplaceAll(search, "'", "''")
	if strings.TrimSpace(search) == "" {
		return "", kurohelpercore.ErrSearchNoContent
	}

	result := "%"
	searchRune := []rune(search)
	for i, r := range searchRune {
		if kurohelpercore.IsEnglish(r) && i < len(searchRune)-1 {
			if kurohelpercore.IsEnglish(searchRune[i+1]) {
				result += string(r)
			} else {
				result += string(r) + "%"
			}
		} else {
			result += string(r) + "%"
		}
	}
	return result, nil
}

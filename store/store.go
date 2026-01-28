package store

import (
	kurohelperdb "kurohelper-db"

	"github.com/sirupsen/logrus"
)

var (
	ZhtwToJp        map[rune]rune
	SeiyaCorrespond map[string]string
)

func InitZhtwToJp() {
	entries, err := kurohelperdb.GetAllZhtwToJp()
	if err != nil {
		logrus.Fatal(err)
	}

	// 轉換
	ZhtwToJp = make(map[rune]rune, len(entries))
	for _, e := range entries {
		keyRunes := []rune(e.ZhTw)
		valRunes := []rune(e.Jp)

		// 確保都是單一字
		if len(keyRunes) == 1 && len(valRunes) == 1 {
			ZhtwToJp[keyRunes[0]] = valRunes[0]
		}
	}
}

func InitSeiyaCorrespond() {
	entries, err := kurohelperdb.GetAllSeiyaCorrespond()
	if err != nil {
		logrus.Fatal(err)
	}

	// 轉換
	SeiyaCorrespond = make(map[string]string, len(entries))
	for _, e := range entries {
		SeiyaCorrespond[e.GameName] = e.SeiyaURL
	}
}

package service

// 因為資料庫問題，暫時不遷移下來，這邊程式碼放置

import (
	"time"

	kurohelperdb "kurohelper-db"

	"kurohelper-core/erogs"

	"gorm.io/gorm"
)

func AddHasPlayedService(erogsData *erogs.FuzzySearchGameResponse, userID string, userName string, completeDate *time.Time) error {
	err := kurohelperdb.Dbs.Transaction(func(tx *gorm.DB) error {
		// 1. 確保 User 存在
		if _, err := kurohelperdb.EnsureUserTx(tx, userID, userName); err != nil {
			return err
		}

		// 2. 確保 Brand 存在
		if _, err := kurohelperdb.EnsureBrandErogsTx(tx, erogsData.BrandID, erogsData.BrandName); err != nil {
			return err
		}

		// 3. 確保 Game 存在
		if _, err := kurohelperdb.EnsureGameErogsTx(tx, erogsData.ID, erogsData.Gamename, erogsData.BrandID); err != nil {
			return err
		}

		// 4. 建立資料
		if err := kurohelperdb.CreateUserHasPlayedTx(tx, userID, erogsData.ID, completeDate); err != nil {
			return err
		}

		return nil // commit
	})
	if err != nil {
		return err
	}
	return nil
}

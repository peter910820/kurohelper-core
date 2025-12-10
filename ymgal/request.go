package ymgal

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	kurohelpercore "github.com/peter910820/kurohelper-core"
	"github.com/sirupsen/logrus"
)

// 做一次重試(取新Token)的版本
func sendWithRetry(apiRoute string) ([]byte, error) {
	r, err := sendGetRequest(apiRoute)
	if err != nil {
		if errors.Is(err, kurohelpercore.ErrYmgalInvalidAccessToken) {
			logrus.Warnf("%s, refreshing and retrying...", err)
			err = GetToken()
			if err != nil {
				return nil, err
			}

			r, err = sendGetRequest(apiRoute)
			if err != nil {
				return nil, err
			}

			return r, nil
		}
		return nil, err
	}
	return r, nil
}

// 統一發送請求(Get版本)
func sendGetRequest(apiRoute string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiRoute, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("version", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, kurohelpercore.ErrYmgalInvalidAccessToken
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w %d", kurohelpercore.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

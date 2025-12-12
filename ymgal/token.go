package ymgal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	kurohelpercore "github.com/peter910820/kurohelper-core"
)

// token
var token tokenResp

// 取得Token
func GetToken() error {
	if !kurohelpercore.IsValidURL(cfg.Endpoint) {
		return ErrCfgInvalid
	}

	req, err := http.NewRequest(http.MethodGet, cfg.Endpoint+fmt.Sprintf("/oauth/token?grant_type=client_credentials&client_id=%s&client_secret=%s&scope=public", cfg.ClientID, cfg.ClientSecret), nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("%w %d", kurohelpercore.ErrStatusCodeAbnormal, resp.StatusCode)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(r, &token)
	if err != nil {
		return err
	}

	return nil
}

package ymgal

import (
	"encoding/json"
	"fmt"
)

type (
	ArchiveResp struct {
		Game       Game       `json:"game"`
		CidMapping CidMapping `json:"cidMapping"`
		PidMapping PidMapping `json:"pidMapping"`
	}

	Game struct {
		PublishVersion int             `json:"publishVersion"`
		PublishTime    string          `json:"publishTime"`
		Publisher      int             `json:"publisher"`
		Name           string          `json:"name"`
		ChineseName    string          `json:"chineseName"`
		ExtensionName  []ExtensionName `json:"extensionName"`
		Introduction   string          `json:"introduction"`
		State          string          `json:"state"`
		Weights        int             `json:"weights"`
		MainImg        string          `json:"mainImg"`
		MoreEntry      any             `json:"moreEntry"`
		Gid            int             `json:"gid"`
		DeveloperID    int             `json:"developerId"`
		HaveChinese    bool            `json:"haveChinese"`
		TypeDesc       string          `json:"typeDesc"`
		ReleaseDate    string          `json:"releaseDate"`
		Restricted     bool            `json:"restricted"`
		Country        string          `json:"country"`
		Website        []Website       `json:"website"`
		Characters     []Character     `json:"characters"`
		Releases       []Release       `json:"releases"`
		Staff          []Staff         `json:"staff"`
		Type           string          `json:"type"`
		Freeze         bool            `json:"freeze"`
	}

	ExtensionName struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Desc string `json:"desc"`
	}

	Website struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	}

	Character struct {
		Cid               int `json:"cid"`
		CvId              int `json:"cvId"`
		CharacterPosition int `json:"characterPosition"`
	}

	Release struct {
		ID               int    `json:"id"`
		ReleaseName      string `json:"releaseName"`
		RelatedLink      string `json:"relatedLink"`
		Platform         string `json:"platform"`
		ReleaseDate      string `json:"releaseDate,omitempty"`
		ReleaseLanguage  string `json:"releaseLanguage,omitempty"`
		RestrictionLevel string `json:"restrictionLevel,omitempty"`
	}

	Staff struct {
		Sid     int    `json:"sid"`
		Pid     int    `json:"pid"`
		EmpName string `json:"empName"`
		EmpDesc string `json:"empDesc"`
		JobName string `json:"jobName"`
	}

	CidMapping struct {
		Cid         int    `json:"cid"`
		Name        string `json:"name"`
		ChineseName string `json:"chineseName,omitempty"`
		MainImg     string `json:"mainImg"`
		State       string `json:"state"`
		Freeze      bool   `json:"freeze"`
	}

	PidMapping struct {
		Pid         int    `json:"pid"`
		Name        string `json:"name"`
		ChineseName string `json:"chineseName,omitempty"`
		MainImg     string `json:"mainImg"`
		State       string `json:"state"`
		Freeze      bool   `json:"freeze"`
	}
)

func Archive(gid int) (*ArchiveResp, error) {
	r, err := sendWithRetry(fmt.Sprintf("%s/open/archive?gid=%d", cfg.Endpoint, gid))
	if err != nil {
		return nil, err
	}

	var result basicResp[ArchiveResp]

	// 如果查不到資料就會直接回傳HTML，無法用status code判斷
	err = json.Unmarshal(r, &result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, ErrAPIFailed{Code: result.Code}
	}

	return &result.Data, nil
}

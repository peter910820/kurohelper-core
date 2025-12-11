package ymgal

type (
	basicResp[T any] struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    T    `json:"data"`
	}

	tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}
)

type (
	randomGameResp struct {
		GID         int    `json:"gid"`
		DeveloperID int    `json:"developerId"`
		Name        string `json:"name"`
		ChineseName string `json:"chineseName"`
		HaveChinese bool   `json:"haveChinese"`
		MainImg     string `json:"mainImg"`
		ReleaseDate string `json:"releaseDate"`
		State       string `json:"state"`
	}
)

type (
	SearchGameResp struct {
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

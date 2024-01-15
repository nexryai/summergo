package summergo

type Player struct {
	Url               string   `json:"url"`
	Width             int      `json:"width"`
	Height            int      `json:"height"`
	IframePermissions []string `json:"allow"`
}

type Summary struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	SiteName    string `json:"sitename"`
	Player      Player `json:"player"`
	Sensitive   bool   `json:"sensitive"`
	ActivityPub string `json:"activitypub"`
}

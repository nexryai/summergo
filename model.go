package summergo

type Player struct {
	Url               string   `json:"url,omitempty"`
	Width             int      `json:"width,omitempty"`
	Height            int      `json:"height,omitempty"`
	IframePermissions []string `json:"allow,omitempty"`
}

type Summary struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Description string `json:"description,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	SiteName    string `json:"sitename"`
	Player      Player `json:"player,omitempty"`
	Sensitive   bool   `json:"sensitive"`
	ActivityPub string `json:"activitypub,omitempty"`
}

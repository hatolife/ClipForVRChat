package appcore

type Mode string

const (
	ModeProcess  Mode = "process"
	ModeSettings Mode = "settings"
	ModeResults  Mode = "results"
	ModeError    Mode = "error"
)

type UIState struct {
	Mode          Mode           `json:"mode"`
	Message       string         `json:"message"`
	ConfigPath    string         `json:"configPath"`
	Config        Config         `json:"config"`
	Results       []Result       `json:"results"`
	History       []HistoryEntry `json:"history"`
	PendingPaths  []string       `json:"pendingPaths"`
	ProcessOnSave bool           `json:"processOnSave"`
}

type Result struct {
	SourcePath       string   `json:"sourcePath"`
	Name             string   `json:"name"`
	OutputPath       string   `json:"outputPath"`
	URL              string   `json:"url"`
	QRURLs           []string `json:"qrUrls"`
	Thumbnail        string   `json:"thumbnail"`
	Error            string   `json:"error"`
	HistoryID        string   `json:"historyId"`
	DiscordMessageID string   `json:"discordMessageId"`
	DiscordWebhookID string   `json:"discordWebhookId"`
	DiscordToken     string   `json:"discordToken"`
}

type ProgressEvent struct {
	Index  int    `json:"index"`
	Total  int    `json:"total"`
	Stage  string `json:"stage"`
	Result Result `json:"result"`
}

type ProcessRequest struct {
	Paths []string `json:"paths"`
}

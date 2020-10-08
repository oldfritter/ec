package models

type Message struct {
	CommonModel
	ThemeId    int    `json:"theme_id"`
	ThemeType  string `json:"theme_type"`
	SenderSn   string `json:"sender_sn"`
	ReceiverSn string `json:"receiver_sn"`
	Content    string `json:"context"`
	Level      int    `json:"level"` // 多层加密，标明层数
}

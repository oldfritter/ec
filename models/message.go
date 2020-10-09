package models

type Message struct {
	ThemeId    int
	ThemeType  string
	SenderSn   string
	ReceiverSn string
	Content    string
	Level      int // 多层加密，标明层数
}

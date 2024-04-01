package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"` // поле optional, поэтомоу указатель
}

type IncomingMessage struct {
	Text string `json:"text"` // команды и ссылки
	From From   `json:"from"` // откуда сообщение
	Chat Chat   `json:"chat"` // куда сообщение
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

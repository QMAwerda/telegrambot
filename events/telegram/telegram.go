package telegram

import (
	"errors"

	"github.com/QMAwerda/telegrambot/clients/telegram"
	"github.com/QMAwerda/telegrambot/events"
	"github.com/QMAwerda/telegrambot/lib/e"
	"github.com/QMAwerda/telegrambot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknowEventType = errors.New("can't process message")
	ErrUnknowMetaType  = errors.New("unknown meta type")
)

// не принято передавать интерфейс по указателю, поэтому storage.Storage
func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{ // offset итак по умолчанию будет 0
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	// возможно, тут стоит вернуть внутреннюю ошибку о том, что updates пустой
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1 // Тогда при след. запросах получим новые сообщения.

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type { // чтобы работать с другими апдейтами тг, нужно просто добавить еще один кейс
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknowEventType)
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process message", err) }()

	meta, err := meta(event)
	if err != nil {
		return err
	}
	// Если пользователь скинул ссылку, ее нужно сохранить, если отправил команду rand.Ind, то нужно найти ссылку и вернуть
	// Если отправит help, нужно дать ему краткую справку по боту, чтоб он понимал, как им пользоваться
	// Назовем все эти действия командами и вынесем в отдельный файл

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta) // проверка, что поле не nil, а именно Meta (приведением его к типу Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknowMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message { // если message, то поля upd.Message точно ненулевые
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res

	// chatID username
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

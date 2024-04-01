package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/QMAwerda/telegrambot/lib/e"
	"github.com/QMAwerda/telegrambot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// это что-то вроде API роутера. Мы будем смотреть на текст сообщения и по его содержанию понимать, какая это команда.
// Эту логику можно вынести в отдельный пакет, но она сильно связана с телеграмом, поэтому оставим ее тут.
func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text) // удаляем пробелы из строки
	// пишем содержимое команды и автора сообщения в лог, чтобы можно было проверить работу.
	log.Printf("got new command '%s' from '%s", text, username)

	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// start: /start: hello + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username) // лучше делать неэкспортируемыми команды, используемые внутри типа
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL, // называем, чтоб не конфликтовать с пакетом url
		UserName: username,
	}

	isExist, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExist)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) { // и при этом есть сохраненные страницы
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) { // нет сохраненных страниц
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page) // удаляем отрправленную ссылку

}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// помним, что ссылки без префикса hhtps:// буду считаьтся невалидными
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

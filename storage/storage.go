package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/QMAwerda/telegrambot/lib/e"
)

// Тут будет интерфейс, он сможет работать и с файловой системой и с базой данных и тд.
// Сейчас будет хранение через файловую систему.

type Storage interface {
	Save(p *Page) error // по ссылке, потому что в теории тип может расширяться
	// можно сделать, чтобы бот переходил по ссылке и делал превью и тд.
	// при передаче по значению, все перечисленные поля будут копироваться
	PickRandom(userName string) (*Page, error) // берет имя, чтоб понимать, чьи ссылки искать
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

// страница, на которую ведет ссылка, которую мы скинули боту
type Page struct {
	URL      string
	UserName string // имя пользователя, который ее скинул
	// Created time.Time - можно добавить поле времени создания
}

// скорее всего, стоит добавить сюда передачу по указателю, для расширяемости в будущем
func (p Page) Hash() (string, error) {
	h := sha1.New()
	// в Функции New создается объект интерфейса Hash, который содержит io.Writer, значит
	// так мы должны передавать в него параметры.
	// Если мы будем хранить все ссылки всех пользователей в одной папке, то может получиться,
	// что у разных пользователей одинаковые ссылки и хеш будет одинаковый. Чтобы этого избежать
	// сделаем хеш по URL + UserName
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	// Конструкция, для перевода байт в строку
	return fmt.Sprintf("%x", h.Sum(nil)), nil // h.Sum() вернет нам байты хеша, nil - т.к. доп байты ему не нужы
}

package sqlite

import (
	"context"
	"database/sql" // common package for work with any sql database
	"fmt"

	"github.com/QMAwerda/telegrambot/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// path to the database
func New(path string) (*Storage, error) {
	// this function return as the err, and some kind of entity(сущность), with witch we will interact with the db
	db, err := sql.Open("sqlite3", path) // specify the database that we will work with
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	// we will try to make the connection with the db, to check is it work or no.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't coonect to the database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	// sql querry, that will save the data to the database
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	// контекст это сущность, которая помогает отменить нам выполнение любой вложенной
	// функции в любой момент времени, когда нам это необходимо.

	// делаем запрос. есть функция проще, без контекста -Exec()
	// использование контекста - хороший тон, он помогает установить таймауты на вложенные вызовы.
	// context, query, arg1, arg2, arg3... агрументы - то, что мы подставляем в запрос вместо знаков вопроса
	// результат работы метода нам не интересен, потому что мы не будем ничего делать с этими данными
	// поэтому ставим прочерк _
	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	// лучше сделать свой рандомайзер, который выдаст рандомное число и по нему сделать OFFSET
	// так будет быстрее работать, чем через ORDER BY RANDOM()
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	// метод QueryRow вернет данные из бд, но в нестандартном формате, потому что бд неизвестно в какой тип мы положим
	// результаты. Чтобы их положить в определенный тип, нужна функция Scan(), в нее передается ссылка на переменную, куда
	// мы хотим положить эти данные

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url) // у нас запись будет ровно одна
	if err == sql.ErrNoRows {                                // Если нет данных
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// возможно стоит удалять страницу по найденному id и добавлять его в бд
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	//Мы не составляем запрос так:
	// DELETE FROM... WHERE url =" + "page.URL" + "..."
	// Потому что злоумышленник может вставить в url что-то типо: addr; DROP DATABASE;
	// И все нам разнесет. Поэтому мы оставляем работу с данными на базу данных, она сама знает, что ей нужно
	// на вход и не будет обрабатывать такие ситуации
	return nil
}

// Ниже будет комментарий в стиле go doc, является хорошим тоном
// Название функции, что делает, в конце точка. Тогда при наведении на функцию увидим короткую справку
// А также можем сгенерировать полноценную документацию с помощью godoc

// IsExist checks if page exist in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	// я скептично отношусь к звездочкам
	q := `SELECT COUNT(*) FROM pages WHERE url = ? and user_name = ?`

	var count int

	err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

// Создаем таблицу в бд
func (s *Storage) Init(ctx context.Context) error {
	// можно создать дополнительно индексы
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table %w", err)
	}

	return nil
}

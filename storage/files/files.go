package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/QMAwerda/telegrambot/lib/e"
	"github.com/QMAwerda/telegrambot/storage"
)

type Storage struct {
	basePath string
}

const defautlPerm = 0774 // permisson - права всем пользователям / full, full, read

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = e.WrapIfErr("can't save page", err)
	}()
	// path.Join() будет корректно ставить слеши на маке и линуксе,
	// но на винде не будет, потому что там "\" в качестве разделителя
	// Для этого придумали пакет filepath. Ниже мы создаем путь до директории, где лежит файл.
	fPath := filepath.Join(s.basePath, page.UserName)       // (путь, имя папки)
	if err := os.MkdirAll(fPath, defautlPerm); err != nil { // создаст все директории, которые входят в переданный путь
		return err
	}
	// чтобы все файлы имели уникальное имя, используем хеш с переданной страницы
	// процедуру получения хеша делаем методом для типа Page

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	// через вызов функции мы получаем возможность поставить конструкцию _=file.Close
	defer func() { _ = file.Close() }() // которая показывает, что мы сознательно игнорируем ошибку

	// Сериализуем страницу - приводим ее к формату, в котором ее можно записать в файл
	// Один из вариантов - перевести структуру в json, а потом его парсить (СТОИТ ДОБАВИТЬ ЭТУ ВОЗМОЖНОСТЬ)
	// А пока работаем с gob - гошный формат файлов, который меньше весит, нужно про него прочитать. Он эффективнее json
	if err := gob.NewEncoder(file).Encode(page); err != nil { // страница будет преобразована в формат gob и записана в указанный файл
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random", err) }()
	// Получаем путь до файлов
	path := filepath.Join(s.basePath, userName)
	// Получаем список файлов
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		// вынесем ошибку в переменную пакета, чтобы ее можно было проверить снаружи
		return nil, storage.ErrNoSavedPages // это потребуется, чтобы бот сказал пользователю, что сохраненных сообщений пока нет
	}

	// Теперь получим случайное число от 0 и до номера последнего файла - 1
	// rand генерирует псевдослучайные числа, ему понадобится Seed(value). Если он всегда будет одинаковый, то ты каждый раз будешь
	// получать одну и ту же последовательность. Чтобы этого не было, будем использовать время
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // определение времени
	n := r.Intn(len(files))                              // rand.Intn(верхняя граница) // rand.Int вернет число от 0 до n-1

	file := files[n]

	// open decode - вынесем в отдельную функцию

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)
	// проверяем существование через функцию os.Stat()
	// она возвращает ошибки, и если ошибка будет как снизу, значит файла нет
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist): // Если файла нет
		return false, nil
	case err != nil: // Если ошибка другая, передаем ее наверх
		msg := fmt.Sprintf("can't check if page %s exists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode a page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page // сюда мы декодируем файл

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

// Чтобы в будущем поменять формирование имен страниц, не нужно будет искать все вызовы функции Hash().
// Будет достаточно изменить ее в этой функции. Например, добавить расширение к файлам.
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

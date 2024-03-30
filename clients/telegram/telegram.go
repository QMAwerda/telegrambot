package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/QMAwerda/telegrambot/lib/e"
)

type Client struct {
	host     string      // хост api сервиса телеграма
	basePath string      // базовый путь, с которого начинаются все запросы
	client   http.Client // храним тут hhtp клиент, чтоб не создавать для каждого запроса отдельно
}

// хост api/путь
// tg-bot.com/bot<token>

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

// Функция создания клиента
func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// Так мы можем генерировать токен в разных местах программы, а если телеграм решит заменить префикс "bot"
// на какой-то другой, то нам не придется рефакторить код во всех местах, где создаем токен
func newBasePath(token string) string {
	return "bot" + token
}

// Экспортирвемые методы принято ставить выше неэкспортируемых в коде
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	// формируем параметры запроса с помощью пакета url
	// для запроса (читаем в документации), нужны limit (сколько апдейтов получить за один запрос) и
	// offset (они хранятся как стек) и если мы уже получили сообщения, то нужно сдвинуться по ним с помощью offset.
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) // метод Add добавляет указанные параметры к запросу
	q.Add("offset", strconv.Itoa(limit))  // мы получаем int в параметрах, но Add() ждет строку, поэтому приводим к ней
	// Itoa = Integer to ASCII

	// Теперь отправим запрос. Код для запроса будет одинаковый для всех методов, поэтому вынесем его в отдельную функцию.

	// doRequest <- getUpdates

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err) // Ошибку вернем только в одном месте, поэтому без defer
		// тут сразу вызываем Wrap, проверка на nil в методе не нужна
	}

	return nil
}

// в ответ получим байты, которые вернет наш запрос
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	// формируем URL на который будет отправляться запрос:
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	// Если делать такую склейку пути: c.basePath + "/" + method, то можно получить и два и три слеша, если в конце basePath
	// или в начале method будет слеш. Чтобы это не проверять, мы используем функцию склейки путей path.Join(...,...). Она
	// уберет лишние слеши и добавит недостающие.

	// Подготавливаем запрос.
	// Первый параметр - сам запрос. Чтобы не ошибиться при написании "GET" используем константу MethodGet,
	// которая содержит эту же строку
	// Второй параметр - url в текстовом виде. Тип url.URL реализует интерфейс Stringer, поэтому у него есть
	// метод String(). Таким образом, вызываем u.String()
	// Третий параметр - тело запроса. Оно обычно отсутствует у метода GET
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// Если возвращать обычную ошибку err, то при записи в лог мы можем даже не понять откуда эта ошибка
		// поэтому мы ее оборачиваем в fmt.Errorf, дописывая дополнительную информацию
		// Посмотри про errors.Is() и errors.As()
		return nil, err
	}
	// Передаем в req параметры запроса, которые мы получили из аргумента (query)
	req.URL.RawQuery = query.Encode() // Он приведет параметры к виду записи query параметров. "a=b&c=e"

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }() // закроем тело ответа, проигнорировав ошибку
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

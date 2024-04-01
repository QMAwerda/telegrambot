package eventconsumer

import (
	"log"
	"time"

	"github.com/QMAwerda/telegrambot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // кол-во событий, которые мы будем обрабатывать за раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

//TODO: можно сделать так, чтобы в случае невыполнения Fetch(), например при проблемах с сетью, у нас делался ретрай
// Первый способ - сделать несколько ретраев с определенным интервалом (ну такое)
// Второй - реализовать ретрай, который будет рости с экспоненциальным временем и, достигнув определенного таймаута, остановится
// Проблема подхода ниже в том, что если итерация идет не раз в секунду, а, например раз в час, то при ошибки в Fetch придется много
// ждать. Либо если код выполняется до метода Fetch, его при continue тоже придется выполнять заново.

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil { // неудачный подход, нужно доработать
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

// TODO: сделать ретрай.
// Проблема метода ниже, в случае ошибки, мы потеряем ссылку, потому что фетчер сдвинет свой offset.
// Тут тоже нужен механизм ретрая, например, стоит сохранять данные во временное хранилище и брать отуда в таком случае.
// Таким образом, хороший вариант, это сделать фолбек в оперативной памяти. Еще способ - подтверждение для Fetcher. Не делать сдвиг
// пока не станет ясно, что все обработано.

// Еще проблема, если ошибок много (например, при проблемах с сетью), то мы можем потратить много времени впустую, возможно, стоит
// сделать счетчик ошибок и если там их будет штук 5, то просто возвращать ошибку, вместо continue.

// Еще, мы обрабатываем все действия один за одним, хотя могли бы сделать их асинхронно.
// Используя примитивы синхронизации, нужно сделать ее в анонимных горутинах, WaitGroup, все дела.

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new events: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event %s", err.Error()) // Чем err.Err() отличается от err.Is() или err.As() ?

			continue
		}
	}

	return nil
}

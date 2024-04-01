package telegram

// вместо кавычек пишем ` ` чтобы сохранять форматирование
const msgHelp = `I can save and keep you pages. Aslo i can offer u them to read.
In order to save page, just send me a link to it with.

In order to get a random page from your list, send me command /rnd.
Caution! After that, this page will be removed from your list!`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command🤔"
	msgNoSavedPages   = "You haven't saved pages😔"
	msgSaved          = "Saved!👌"
	msgAlreadyExist   = "You already have this page in your list😴"
)

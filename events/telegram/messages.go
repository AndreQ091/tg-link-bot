package telegram

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have not saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"

	msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me al link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After that, this page will be removed from your list!`
	msgHello = "Hi there! 👾\n\n" + msgHelp
)

package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	e "tg-link-bot/lib/error"
	"tg-link-bot/storage"
)

const (
	RndCMD   = "/rnd"
	HelpCMD  = "/help"
	StartCMD = "/start"
)

func (p *Processor) doCMD(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from %s", text, username)

	if isAddCMD(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCMD:
		return p.sendRandom(chatID, username)
	case HelpCMD:
		return p.sendHelp(chatID)
	case StartCMD:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)

	}

}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("can't do command: save page", err)
	}()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)

	if err != nil {
		return err
	}

	if isExists {
		p.tg.SendMessage(chatID, msgAlreadyExists)
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
	defer func() {
		err = e.WrapIfErr("can't send random page", err)
	}()

	page, err := p.storage.PickRandom(username)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCMD(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

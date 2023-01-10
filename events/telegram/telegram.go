package telegram

import (
	"errors"
	"tg-link-bot/clients/telegram"
	"tg-link-bot/events"
	e "tg-link-bot/lib/error"
	"tg-link-bot/storage"
)

type Processor struct {
	tg      *tgClient.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrorUknownEventType = errors.New("uknown event type")
	ErrorUknownMetaType  = errors.New("uknown meta type")
)

func NewProcessor(client *tgClient.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)

	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrorUknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCMD(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrorUknownMetaType)
	}

	return res, nil
}

func event(update tgClient.Update) events.Event {
	updType := fetchType(update)

	res := events.Event{
		Type: updType,
		Text: fetchText(update),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	}

	return res
}

func fetchType(update tgClient.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update tgClient.Update) string {
	if update.Message == nil {
		return " "
	}
	return update.Message.Text
}

package events

import (
	"context"
	"log/slog"

	"github.com/vovanwin/template/internal/pkg/centrifugo"
)

type Event struct {
	UserID  string `json:"-"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Bus struct {
	client *centrifugo.Client
	log    *slog.Logger
}

func NewBus(client *centrifugo.Client, log *slog.Logger) *Bus {
	return &Bus{
		client: client,
		log:    log,
	}
}

func (b *Bus) Publish(e Event) error {
	channel := centrifugo.PersonalChannel(e.UserID)
	if err := b.client.Publish(context.Background(), channel, e); err != nil {
		b.log.Warn("failed to publish event", slog.String("userID", e.UserID), slog.Any("err", err))
		return err
	}
	return nil
}

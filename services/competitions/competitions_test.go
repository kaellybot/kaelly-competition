package competitions

import (
	"testing"

	amqp "github.com/kaellybot/kaelly-amqp"
)

// stubMapService records whether the handler ran, so a test can assert that an
// unsupported game never reaches the KTArena map logic.
type stubMapService struct {
	called *int
}

func (stub stubMapService) GetMapRequest(_ amqp.Context, _ *amqp.RabbitMQMessage) {
	*stub.called++
}

func newTestService(called *int, replies *[]*amqp.RabbitMQMessage) *Impl {
	broker := amqp.Mock{
		ReplyFunc: func(msg *amqp.RabbitMQMessage, _, _ string) error {
			*replies = append(*replies, msg)
			return nil
		},
	}
	return New(&broker, stubMapService{called: called})
}

func TestConsumeHandlesDofus(t *testing.T) {
	t.Parallel()

	called := 0
	replies := make([]*amqp.RabbitMQMessage, 0)
	service := newTestService(&called, &replies)

	service.consume(amqp.Context{CorrelationID: "correlationID"}, &amqp.RabbitMQMessage{
		Type: amqp.RabbitMQMessage_COMPETITION_MAP_REQUEST,
		Game: amqp.Game_DOFUS_GAME,
	})

	if called != 1 {
		t.Errorf("handler called %d time(s), want 1", called)
	}
	if len(replies) != 0 {
		t.Errorf("consume published %d reply(ies), want 0", len(replies))
	}
}

// KTArena maps are Dofus-only: Touch must get a FAILED reply tagged DOFUS_TOUCH,
// never a Dofus map.
func TestConsumeRefusesUnsupportedGames(t *testing.T) {
	t.Parallel()

	for _, game := range []amqp.Game{amqp.Game_DOFUS_TOUCH, amqp.Game_DOFUS_RETRO} {
		t.Run(game.String(), func(t *testing.T) {
			t.Parallel()

			called := 0
			replies := make([]*amqp.RabbitMQMessage, 0)
			service := newTestService(&called, &replies)

			service.consume(amqp.Context{CorrelationID: "correlationID"}, &amqp.RabbitMQMessage{
				Type: amqp.RabbitMQMessage_COMPETITION_MAP_REQUEST,
				Game: game,
			})

			if called != 0 {
				t.Errorf("handler called %d time(s), want 0", called)
			}
			if len(replies) != 1 {
				t.Fatalf("consume published %d reply(ies), want 1", len(replies))
			}
			if replies[0].GetStatus() != amqp.RabbitMQMessage_FAILED {
				t.Errorf("Status = %v, want %v", replies[0].GetStatus(), amqp.RabbitMQMessage_FAILED)
			}
			if replies[0].GetGame() != game {
				t.Errorf("Game = %v, want %v", replies[0].GetGame(), game)
			}
			if replies[0].GetType() != amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER {
				t.Errorf("Type = %v, want %v", replies[0].GetType(),
					amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER)
			}
		})
	}
}

func TestConsumeDropsRequestWithoutGame(t *testing.T) {
	t.Parallel()

	called := 0
	replies := make([]*amqp.RabbitMQMessage, 0)
	service := newTestService(&called, &replies)

	service.consume(amqp.Context{CorrelationID: "correlationID"},
		&amqp.RabbitMQMessage{Type: amqp.RabbitMQMessage_COMPETITION_MAP_REQUEST})

	if called != 0 {
		t.Errorf("handler called %d time(s), want 0", called)
	}
	if len(replies) != 0 {
		t.Errorf("consume published %d reply(ies), want 0", len(replies))
	}
}

func TestConsumeIgnoresUnknownType(t *testing.T) {
	t.Parallel()

	called := 0
	replies := make([]*amqp.RabbitMQMessage, 0)
	service := newTestService(&called, &replies)

	service.consume(amqp.Context{CorrelationID: "correlationID"}, &amqp.RabbitMQMessage{
		Type: amqp.RabbitMQMessage_ABOUT_REQUEST,
		Game: amqp.Game_DOFUS_GAME,
	})

	if called != 0 {
		t.Errorf("handler called %d time(s), want 0", called)
	}
	if len(replies) != 0 {
		t.Errorf("consume published %d reply(ies), want 0", len(replies))
	}
}

func TestIsGameSupported(t *testing.T) {
	t.Parallel()

	tests := map[amqp.Game]bool{
		amqp.Game_ANY_GAME:    false,
		amqp.Game_DOFUS_GAME:  true,
		amqp.Game_DOFUS_TOUCH: false,
		amqp.Game_DOFUS_RETRO: false,
	}

	for game, want := range tests {
		if got := isGameSupported(game); got != want {
			t.Errorf("isGameSupported(%v) = %v, want %v", game, got, want)
		}
	}
}

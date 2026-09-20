package mappers

import (
	"testing"

	amqp "github.com/kaellybot/kaelly-amqp"
)

// The answer must carry the request's game even though competition only serves
// Dofus today: the reply is routed back by game, not by the service's identity.
func TestMapGetMapAnswerCarriesRequestGame(t *testing.T) {
	t.Parallel()

	for _, game := range []amqp.Game{amqp.Game_DOFUS_GAME, amqp.Game_DOFUS_TOUCH} {
		t.Run(game.String(), func(t *testing.T) {
			t.Parallel()

			request := &amqp.RabbitMQMessage{
				Type:                  amqp.RabbitMQMessage_COMPETITION_MAP_REQUEST,
				Game:                  game,
				Language:              amqp.Language_FR,
				CompetitionMapRequest: &amqp.CompetitionMapRequest{MapNumber: 12},
			}

			answer := MapGetMapAnswer(request, 12)

			if answer.GetGame() != game {
				t.Errorf("Game = %v, want %v", answer.GetGame(), game)
			}
			if !answer.IsGameSet() {
				t.Error("answer game must be set, otherwise the broker refuses to publish it")
			}
			if answer.GetLanguage() != amqp.Language_FR {
				t.Errorf("Language = %v, want %v", answer.GetLanguage(), amqp.Language_FR)
			}
			if answer.GetStatus() != amqp.RabbitMQMessage_SUCCESS {
				t.Errorf("Status = %v, want %v", answer.GetStatus(), amqp.RabbitMQMessage_SUCCESS)
			}
			if answer.GetCompetitionMapAnswer().GetMapNumber() != 12 {
				t.Errorf("MapNumber = %v, want 12", answer.GetCompetitionMapAnswer().GetMapNumber())
			}
		})
	}
}

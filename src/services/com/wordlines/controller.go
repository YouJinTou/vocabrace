package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/YouJinTou/vocabrace/games/wordlines"
	"github.com/YouJinTou/vocabrace/services/com/state/data"
	"github.com/YouJinTou/vocabrace/services/com/ws"
	"github.com/google/uuid"
)

// Controller communicaes data between the game logic object and the communication layer.
type Controller struct{}
type cell struct {
	CellIndex int    `json:"c"`
	TileID    string `json:"t"`
	Blank     string `json:"b"`
}
type player struct {
	ID     string `json:"i"`
	Name   string `json:"n"`
	Points int    `json:"p"`
	Tiles  int    `json:"t"`
}
type turn struct {
	IsPlace       bool     `json:"p"`
	IsExchange    bool     `json:"e"`
	IsPass        bool     `json:"q"`
	Word          []*cell  `json:"w"`
	ExchangeTiles []string `json:"t"`
}
type result struct {
	g   *wordlines.Game
	d   wordlines.DeltaState
	err error
}

// OnStart executes logic at the start of the game.
func (s Controller) OnStart(i data.OnStartInput) data.OnStartOutput {
	c := i.Connections
	players := s.loadPlayers(c)
	game := wordlines.NewSpiralGame(c.Language(), players, wordlines.NewDynamoValidator())
	projected := s.setPlayerData(game.Players)
	messages := []*ws.Message{}
	startState := struct {
		Tiles          *wordlines.Tiles `json:"t"`
		ToMoveID       string           `json:"m"`
		Players        []*player        `json:"p"`
		YourMove       bool             `json:"y"`
		PoolID         string           `json:"pid"`
		TilesRemaining int              `json:"r"`
		Language       string           `json:"z"`
		IsStart        bool             `json:"s"`
	}{nil, game.ToMoveID, projected, false, uuid.New().String(), game.Bag.Count(), c.Language(), true}

	for _, p := range game.Players {
		startState.Tiles = p.Tiles
		startState.YourMove = game.ToMoveID == p.ID
		messages = append(messages,
			ws.NewMessage(c.Domain(), *c.IDByUserID(p.ID), startState, &p.ID))
	}

	return data.OnStartOutput{
		PoolID:   startState.PoolID,
		Messages: messages,
		Game:     game,
	}
}

// OnAction executes logic at each turn.
func (s Controller) OnAction(i data.OnActionInput) data.OnActionOutput {
	turn := turn{}
	bErr := json.Unmarshal([]byte(i.Body), &turn)

	if bErr != nil {
		return data.OnActionOutput{Error: s.error(i, "invalid payload", bErr)}
	}

	game := &wordlines.Game{}
	dynamodbattribute.UnmarshalMap(i.State, game)
	game.SetValidator(wordlines.NewDynamoValidator())
	game.SetLayout("spiral")

	if vErr := s.validateTurn(i, game); vErr != nil {
		return data.OnActionOutput{Error: s.error(i, "invalid turn", vErr)}
	}

	var r *result
	if turn.IsExchange {
		r = s.exchange(&turn, game)
	} else if turn.IsPass {
		r = s.pass(game)
	} else if turn.IsPlace {
		r = s.place(&turn, game)
	} else {
		r = &result{nil, wordlines.DeltaState{}, errors.New("invalid action")}
	}

	if r.err != nil {
		return data.OnActionOutput{Error: s.error(i, "invalid move", r.err)}
	}

	messages := []*ws.Message{ws.NewMessage(
		i.Connections.Domain(),
		i.Initiator,
		r.d.WithPersonal(),
		i.Connections.UserIDByID(i.Initiator)),
	}
	for _, cid := range i.Connections.OtherIDs(i.Initiator) {
		r.d.YourMove = game.ToMoveID == *i.Connections.UserIDByID(cid)
		messages = append(messages, ws.NewMessage(
			i.Connections.Domain(),
			cid,
			r.d.WithoutPersonal(),
			i.Connections.UserIDByID(cid)))
	}

	return data.OnActionOutput{
		Messages: messages,
		Game:     game,
		IsOver:   game.IsOver(),
	}
}

func (s Controller) exchange(turn *turn, g *wordlines.Game) *result {
	game, err := g.Exchange(turn.ExchangeTiles)
	if err != nil {
		return &result{&game, wordlines.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s Controller) pass(g *wordlines.Game) *result {
	game := g.Pass()
	return &result{&game, game.GetDelta(), nil}
}

func (s Controller) place(turn *turn, g *wordlines.Game) *result {
	cells := []*wordlines.Cell{}
	for _, c := range turn.Word {
		cells = append(cells, wordlines.NewCell(
			wordlines.NewTileWithID(c.TileID, c.Blank, 0),
			c.CellIndex,
		))
	}
	game, err := g.Place(wordlines.NewWord(cells))
	if err != nil {
		return &result{&game, wordlines.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s Controller) loadPlayers(c *data.Connections) []*wordlines.Player {
	result := []*wordlines.Player{}
	for _, u := range c.UserIDs() {
		result = append(result, &wordlines.Player{
			ID:     u,
			Name:   u,
			Points: 0,
		})
	}
	return result
}

func (s Controller) setPlayerData(players []*wordlines.Player) []*player {
	result := []*player{}
	for _, p := range players {
		result = append(result, &player{
			ID:     p.ID,
			Name:   p.Name,
			Points: p.Points,
			Tiles:  p.Tiles.Count(),
		})
	}
	return result
}

func (s Controller) validateTurn(data data.OnActionInput, g *wordlines.Game) error {
	if data.InitiatorUserID != g.ToMoveID {
		return errors.New("invalid player turn")
	}
	return nil
}

func (s Controller) error(data data.OnActionInput, message string, err error) *ws.Message {
	return ws.NewErrorMessage(
		data.Connections.Domain(),
		data.Initiator,
		fmt.Sprintf("%s: %s", message, err.Error()),
		data.Connections.UserIDByID(data.Initiator),
	)
}

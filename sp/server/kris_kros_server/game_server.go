package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/game"
	"ups/sp/server/protocol/def"
	"ups/sp/server/protocol/impl"
	"ups/sp/server/protocol/messages"
	"ups/sp/server/protocol/responses"
)

type GameServer struct {
	gamesByLobbyID  map[int]*game.Game
	gamesByPlayerID map[int]*game.Game
	server          *KrisKrosServer
}

func NewGameServer(server *KrisKrosServer) *GameServer {
	return &GameServer{
		server:          server,
		gamesByLobbyID:  make(map[int]*game.Game),
		gamesByPlayerID: make(map[int]*game.Game),
	}
}

func (server *GameServer) CreateGame(players []game.Player) {
	log.Infoln("Starting a game...")

	game := game.Game{
		Players: players,
	}
	err := game.Start()

	if err != nil {
		log.Errorf("Could not start game because of error: '%s'", err)
		return
	}

	for _, player := range players {
		server.gamesByPlayerID[player.ID] = &game
		if player.ID != game.CurrentPlayer.ID {
			server.server.Router.UserStates[player.ID] = PlayerWaitingState{}
		} else {
			server.server.Router.UserStates[player.ID] = PlayersTurnState{}
			server.server.Router.IgnoreTransitionStateChange = true
		}
	}

	for id, letterBag := range game.PlayerIdToPlayerBag {
		log.Infof("Player %d letterbag:", id)
		for _, letter := range letterBag {
			log.Infof("%s of %d points", letter.Value, letter.Points)
		}
		log.Infoln()
		resp := impl.StructMessageResponse(responses.GameStartedResponse{
			Players:        players,
			Letters:        letterBag,
			ActivePlayerID: game.CurrentPlayer.ID,
		})

		server.server.Send(resp, id, 0)
	}
}

func (server *GameServer) OnLetterPlaced(userId int, message messages.LetterPlacedMessage) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	if g.CurrentPlayer.ID != userId {
		message := fmt.Sprintf("It is not turn of player of ID %d ", userId)
		log.Errorln(message)
		return impl.ErrorResponse(message, impl.NotPlayersTurn)
	}

	err := g.HandleSetAtEvent(game.SetLetterAtEvent{
		PlayerID: userId,
		Row:      message.Row,
		Column:   message.Column,
		Letter:   message.Letter.Value,
	})

	if err == nil {

		updatedTiles := []game.Tile{}

		for tile, _ := range g.Desk.CurrentLetters.List {
			updatedTiles = append(updatedTiles, tile)
		}

		for _, player := range g.Players {
			server.server.Send(impl.StructMessageResponse(responses.TilesUpdatedResponse{Tiles: updatedTiles}), player.ID, 0)
			//if player.ID != g.CurrentPlayer.ID {
			//	server.server.Send(impl.StructMessageResponse(responses.TileUpdatedResponse{Tile: g.Desk.Tiles[message.Row][message.Column]}), player.ID, 1)
			//}
		}
		return impl.SuccessResponse("Placed successfully")
	} else {
		log.Errorf("Error while setting a letter: '%s'", err)
		return impl.ErrorResponse(err.Error(), 999)
	}
}

func (server *GameServer) OnLetterRemoved(userId int, message messages.LetterRemovedMessage) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	if g.CurrentPlayer.ID != userId {
		message := fmt.Sprintf("It is not turn of player of ID %d ", userId)
		log.Errorln(message)
		return impl.ErrorResponse(message, impl.NotPlayersTurn)
	}

	err := g.HandleResetAtEvent(game.ResetAtEvent{
		PlayerID: userId,
		Row:      message.Row,
		Column:   message.Column,
	})

	if err == nil {

		g.Desk.ClearCurrentWords()

		for tile, _ := range g.Desk.PlacedLetter.List {
			g.Desk.GetWordsAt(tile.Row, tile.Column)
		}

		updatedTiles := []game.Tile{}

		for tile, _ := range g.Desk.CurrentLetters.List {
			updatedTiles = append(updatedTiles, tile)
		}

		updatedTiles = append(updatedTiles, g.Desk.Tiles[message.Row][message.Column])

		for _, player := range g.Players {
			if player.ID != g.CurrentPlayer.ID {
				server.server.Send(impl.StructMessageResponse(responses.TilesUpdatedResponse{Tiles: updatedTiles}), player.ID, 1)
			}
		}
		return impl.SuccessResponse("Letter removed successfully")
	} else {
		log.Errorf("Error while setting a letter: '%s'", err)
		return impl.ErrorResponse(err.Error(), 999)
	}
}

func (server *GameServer) OnFinishRound(userId int) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	if g.CurrentPlayer.ID != userId {
		message := fmt.Sprintf("It is not turn of player of ID %d ", userId)
		log.Errorln(message)
		return impl.ErrorResponse(message, impl.NotPlayersTurn)
	}

	if len(g.Desk.PlacedLetter.List) <= 0 {
		g.EmptyRounds++
	}

	for _, player := range g.Players {
		if player.ID != userId {
			server.server.Send(impl.StructMessageResponse(responses.PlayerFinishedResponse{}), player.ID, 0)
		}
	}

	return impl.SuccessResponse("Finished successfully")
}

func (server *GameServer) OnApproveWords(userId int) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	player, exists := g.PlayersMap[userId]

	if player.ID == g.CurrentPlayer.ID {
		return impl.ErrorResponse("The player who's the current round cannot accept his own words.", impl.PlayerCannotAcceptHisOwnWords)
	}

	if exists {
		roundAccepted := g.AcceptTurn(player)

		if !roundAccepted {
			for _, player := range g.Players {
				if player.ID != userId {
					server.server.Send(impl.StructMessageResponse(responses.PlayerAcceptedRoundResponse{PlayerID: userId}), player.ID, 0)
				}
			}
		} else {
			g.Next()

		}

		return impl.SuccessResponse("Successfully accepted words...")
	}
	return impl.ErrorResponse(fmt.Sprint("Could not find a player of ID %d", userId), impl.PlayerNotFound)
}

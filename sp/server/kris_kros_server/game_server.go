package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/game"
	"ups/sp/server/model"
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

// creates a new game server
func NewGameServer(server *KrisKrosServer) *GameServer {
	return &GameServer{
		server:          server,
		gamesByLobbyID:  make(map[int]*game.Game),
		gamesByPlayerID: make(map[int]*game.Game),
	}
}

// creates a new game
func (server *GameServer) CreateGame(players []game.Player) {
	log.Infoln("Creating a new game")

	g := game.Game{
		Players: players,
	}

	err := g.Start()

	if err != nil {
		log.Errorf("Could not start game because of error: '%s'", err)
		return
	}

	log.Debugf("Current player on game start name ID=%d, name=%s", g.CurrentPlayer.ID, g.CurrentPlayer.Name)

	for _, player := range players {
		server.gamesByPlayerID[player.ID] = &g
		if player.ID != g.CurrentPlayer.ID {
			server.server.Router.UserStates[player.ID] = PlayerWaitingState{}
		} else {
			server.server.Router.UserStates[player.ID] = PlayersTurnState{}
			server.server.Router.IgnoreTransitionStateChange = true
		}
	}

	// notify players that the game ha started
	for id, letterBag := range g.PlayerIdToPlayerBag {
		log.Debugf("Player's ID=%d letter bag:", id)
		for _, letter := range letterBag {
			log.Debugf("%s of %d points", letter.Value, letter.Points)
		}
		resp := impl.StructMessageResponse(responses.GameStartedResponse{
			Players:        players,
			Letters:        letterBag,
			ActivePlayerID: g.CurrentPlayer.ID,
		})

		server.server.Send(resp, id, 0)
	}
}

// handler of letter placed events
func (server *GameServer) OnLetterPlaced(userId int, message messages.LetterPlacedMessage) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	log.Infof("OnLetterPlaced event by player ID=%d", userId)

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

		// create a list of tiles that were updated
		for tile, _ := range g.Desk.CurrentLetters.List {
			updatedTiles = append(updatedTiles, tile)
		}

		points := g.Desk.GetTotalPoints()
		playerTotalPoints := g.PointsTable[g.CurrentPlayer.ID] + points

		// notify players about tiles that were updated
		for _, player := range g.Players {
			server.server.Send(impl.StructMessageResponse(responses.TilesUpdatedResponse{Tiles: updatedTiles, CurrentPlayerTotalPoints: playerTotalPoints, CurrentPlayerPoints: points}), player.ID, 0)
		}
		return impl.SuccessResponse("Placed successfully")
	} else {
		log.Errorf("Error while setting a letter: '%s'", err)
		return impl.ErrorResponse(err.Error(), impl.LetterCannotBePlaced)
	}
}

// handler of letter removed events
func (server *GameServer) OnLetterRemoved(userId int, message messages.LetterRemovedMessage) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player of ID=%d", userId), impl.GameNotFoundByPlayerId)
	}

	log.Infof("OnLetterPlaced event by player ID=%d", userId)

	if g.CurrentPlayer.ID != userId {
		message := fmt.Sprintf("It is not turn of player of ID=%d ", userId)
		log.Errorln(message)
		return impl.ErrorResponse(message, impl.NotPlayersTurn)
	}

	err := g.HandleResetAtEvent(game.ResetAtEvent{
		PlayerID: userId,
		Row:      message.Row,
		Column:   message.Column,
	})

	if err == nil {
		// clear current letters so we know which words may have been lost in the process of removing a letter
		g.Desk.ClearCurrentLetters()

		for tile, _ := range g.Desk.PlacedLetter.List {
			g.Desk.GetWordsAt(tile.Row, tile.Column)
		}

		updatedTiles := []game.Tile{}

		for tile, _ := range g.Desk.CurrentLetters.List {
			updatedTiles = append(updatedTiles, tile)
		}

		// the tile from which the letter was remove was also updated so we append it
		updatedTiles = append(updatedTiles, g.Desk.Tiles[message.Row][message.Column])

		points := g.Desk.GetTotalPoints()
		playerTotalPoints := g.PointsTable[g.CurrentPlayer.ID] + points

		for _, player := range g.Players {
			server.server.Send(impl.StructMessageResponse(responses.TilesUpdatedResponse{Tiles: updatedTiles, CurrentPlayerPoints: points, CurrentPlayerTotalPoints: playerTotalPoints}), player.ID, 1)
		}
		return impl.SuccessResponse("Letter removed successfully")
	} else {
		log.Errorf("Error while setting a letter: '%s'", err)
		return impl.ErrorResponse(err.Error(), 999)
	}
}

// handler of finish round events
func (server *GameServer) OnFinishRound(userId int) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID=%d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	log.Infof("OnFinishRound event by player ID=%d", userId)

	if g.CurrentPlayer.ID != userId {
		message := fmt.Sprintf("It is not turn of player ID=%d ", userId)
		log.Errorln(message)
		return impl.ErrorResponse(message, impl.NotPlayersTurn)
	}

	if len(g.Desk.PlacedLetter.List) <= 0 {
		// when no letters were placed, end the round
		g.EmptyRounds++
		// when all players failed to place any letters or the last player disconnected, end the game
		if g.EmptyRounds >= g.ActivePlayerCount() && !(g.CurrentPlayer.Disconnected && g.ActivePlayerCount() > 0) {
			log.Infof("Ending a game by a FinishRoundEvent")
			return server.EndGame(g, userId)
		} else {
			server.NextRound(g)
			log.Debugf("Skipped word acceptance, placed letter list is empty, and the current player is ID=%d name=%s", g.CurrentPlayer.ID, g.CurrentPlayer.Name)
			return impl.StructMessageResponse(responses.NewRoundResponse{ActivePlayerID: g.CurrentPlayer.ID})
		}
	} else {
		g.RoundFinished = true
		g.EmptyRounds = 0
		if (g.CurrentPlayer.Disconnected && g.ActivePlayerCount() > 0) || g.ActivePlayerCount() > 1 {
			// when there are other players, notify them that they should decide words validity
			for _, player := range g.Players {
				if player.ID != userId {
					server.server.Router.UserStates[player.ID] = ApproveWordsState{}
					server.server.Send(impl.StructMessageResponse(responses.RoundFinishedResponse{}), player.ID, 0)
				}
			}
		} else {
			// one player plays alone, no one is there to decide validity of his words
			server.NextRound(g)
			return impl.StructMessageResponse(responses.NewRoundResponse{ActivePlayerID: g.CurrentPlayer.ID})
		}
	}

	return impl.SuccessResponse("Finished successfully")
}

// ends the game by the given player id
func (server *GameServer) EndGame(g *game.Game, userId int) def.Response {
	log.Infof("EndGame event by player ID=%d", userId)

	pointsToPlayerMap := make(map[int]game.Player)
	for _, player := range g.Players {
		pointsToPlayerMap[g.PointsTable[player.ID]] = player
	}
	resp := impl.StructMessageResponse(responses.GameEndedResponse{PlayerPoints: pointsToPlayerMap})
	server.server.Router.IgnoreTransitionStateChange = true
	for _, player := range g.Players {
		delete(server.gamesByPlayerID, player.ID)
		server.server.Router.UserStates[player.ID] = AuthorizedState{}
		if player.ID != userId {
			server.server.Send(resp, player.ID, 0)
		}
	}
	return resp
}

// on approve words event handler
func (server *GameServer) OnApproveWords(userId int) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	log.Infof("OnApproveWords event by player ID=%d", userId)

	player, exists := g.PlayersMap[userId]

	if player.ID == g.CurrentPlayer.ID {
		return impl.ErrorResponse("The player who's the current round cannot accept his own words.", impl.PlayerCannotAcceptHisOwnWords)
	}

	if exists {
		roundAccepted := g.AcceptTurn(player)

		if !roundAccepted {
			// notify other players that the words got accepted by this player
			for _, player := range g.Players {
				if player.ID != userId {
					server.server.Send(impl.StructMessageResponse(responses.PlayerAcceptedRoundResponse{PlayerID: userId}), player.ID, 0)
				}
			}
		} else {
			// if the accept of words resulted in a acceptance of around, next round is
			server.NextRound(g)
			return impl.StructMessageResponse(responses.AcceptResultedInNewRound{})
		}

		return impl.SuccessResponse("Successfully accepted words...")
	}
	return impl.ErrorResponse(fmt.Sprintf("Could not find a player of ID %d", userId), impl.PlayerNotFound)
}

// starts a new round of the given game
func (server *GameServer) NextRound(g *game.Game) {
	if g.ActivePlayerCount() > 0 {
		log.Infof("NextRound, current player ID=%d", g.CurrentPlayer.ID)
		g.Next()
		server.server.Send(impl.StructMessageResponse(responses.YourNewRoundResponse{Letters: g.PlayerIdToPlayerBag[g.CurrentPlayer.ID]}), g.CurrentPlayer.ID, 0)
		server.server.Router.IgnoreTransitionStateChange = true
		for _, player := range g.Players {
			if player.ID != g.CurrentPlayer.ID {
				server.server.Send(impl.StructMessageResponse(responses.NewRoundResponse{ActivePlayerID: g.CurrentPlayer.ID}), player.ID, 0)
				server.server.Router.UserStates[player.ID] = PlayerWaitingState{}
			} else {
				server.server.Router.UserStates[player.ID] = PlayersTurnState{}
			}
		}
	}
}

// on player declined words
func (server *GameServer) OnDeclineWords(userId int) def.Response {
	g, exists := server.gamesByPlayerID[userId]

	if !exists {
		log.Errorf("Could not find a game by player ID of %d", userId)
		return impl.ErrorResponse(fmt.Sprintf("Could not find a game by player ID of %d", userId), impl.GameNotFoundByPlayerId)
	}

	log.Infof("OnDeclineWords event by player ID=%d", userId)

	playerThatDeclined, exists := g.PlayersMap[userId]

	if playerThatDeclined.ID == g.CurrentPlayer.ID {
		return impl.ErrorResponse("The player who's the current round cannot decline his own words.", impl.PlayerCannotAcceptHisOwnWords)
	}

	if exists {
		g.WordsDeclined()

		if g.CurrentPlayer.Disconnected {
			// if the current player disconnected all his letters are removed because he cannot edit his letters after denial
			updatedTiles := []game.Tile{}
			for tile, _ := range g.Desk.PlacedLetter.List {
				g.Desk.Tiles[tile.Row][tile.Column].Set = false
			}
			for tile, _ := range g.Desk.CurrentLetters.List {
				g.Desk.Tiles[tile.Row][tile.Column].Highlighted = false
				updatedTiles = append(updatedTiles, g.Desk.Tiles[tile.Row][tile.Column])
			}

			for _, player := range g.Players {
				if player.ID != g.CurrentPlayer.ID {
					server.server.Send(impl.StructMessageResponse(responses.TilesUpdatedResponse{
						Tiles:                    updatedTiles,
						CurrentPlayerPoints:      0,
						CurrentPlayerTotalPoints: g.PointsTable[g.CurrentPlayer.ID],
					}), player.ID, 0)
				}
			}
			g.Desk.CurrentLetters.Clear()
			g.Desk.PlacedLetter.Clear()
			server.NextRound(g)
		} else {
			// users are notified and the current player has to edit his letters
			for _, player := range g.Players {
				if player.ID != userId {
					server.server.Send(impl.StructMessageResponse(responses.PlayerDeclinedWordsResponse{
						PlayerID:   userId,
						PlayerName: playerThatDeclined.Name,
					}), player.ID, 0)
				}
				if player.ID != g.CurrentPlayer.ID {
					server.server.Router.UserStates[player.ID] = PlayerWaitingState{}
				} else {
					server.server.Router.UserStates[player.ID] = PlayersTurnState{}
				}
			}
		}
		return impl.SuccessResponse("Successfully declined words...")
	}
	return impl.ErrorResponse(fmt.Sprintf("Could not find a player of ID %d", userId), impl.PlayerNotFound)
}

// a player left the game
func (server *GameServer) PlayerLeft(playerID int, stateID int, playerLeaving bool) {
	log.Infof("PlayerLeft event by player ID=%d, leaving=%s", playerID, playerLeaving)
	g, exists := server.gamesByPlayerID[playerID]
	if exists {
		_, exists := g.PlayersMap[playerID]
		if exists {
			// notify players that a player left
			for index, player := range g.Players {
				if player.ID == playerID {
					log.Warnf("Marking playerID=%d playerName=%s as disconnected", player.ID, player.Name)
					g.Players[index].Disconnected = true
					g.PlayersMap[playerID] = g.Players[index]
					if g.CurrentPlayer.ID == playerID {
						g.CurrentPlayer = g.Players[index]
					}
					break
				}
			}
			log.Warnf("active player count => %d", g.ActivePlayerCount())
			// decide what to do based on the player's state
			if g.ActivePlayerCount() > 0 {
				switch stateID {

				case PLAYERS_TURN_STATE_ID:
					server.OnFinishRound(playerID)

				case APPROVE_WORDS_STATE_ID:
					server.OnApproveWords(playerID)

				}
			} else {
				// if he was the last player, end the game
				server.EndGame(g, -1)
			}
		}
	}

	if playerLeaving {
		// remove the player if he left on his own will
		leavingPlayerIndex := -1
		for index, player := range g.Players {
			if playerID == player.ID {
				leavingPlayerIndex = index
				break
			}
		}
		g.Players = append(g.Players[:leavingPlayerIndex], g.Players[leavingPlayerIndex+1:]...)
	}

	for _, player := range g.Players {
		if !player.Disconnected {
			server.server.Send(impl.StructMessageResponse(responses.PlayerConnectionChanged{PlayerID: playerID, Disconnected: true}), player.ID, 0)
		}
	}
}

// on player reconnected
func (server *GameServer) PlayerReconnected(playerID int) def.Response {
	log.Infof("PlayerReconnected event by player ID=%d", playerID)
	g, exists := server.gamesByPlayerID[playerID]

	if exists {
		var resp def.Response
		for index, player := range g.Players {
			if player.ID == playerID {
				g.Players[index].Disconnected = false
				g.PlayersMap[playerID] = g.Players[index]

				tiles := []game.Tile{}

				for row := 0; row < game.DeskSize; row++ {
					for column := 0; column < game.DeskSize; column++ {
						if g.Desk.Tiles[row][column].Set {
							tiles = append(tiles, g.Desk.Tiles[row][column])
						}
					}
				}

				pointsToPlayerMap := make(map[int]game.Player)

				for _, player := range g.Players {
					pointsToPlayerMap[g.PointsTable[player.ID]] = player
				}

				playerIDsThatAccepted := []int{}

				for player, _ := range g.PlayersThatAccepted.List {
					playerIDsThatAccepted = append(playerIDsThatAccepted, player.ID)
				}

				// send him the current state of the game
				resp = impl.StructMessageResponse(responses.GameStateRegenerationResponse{
					Tiles:                 tiles,
					ActivePlayerID:        g.CurrentPlayer.ID,
					PlayerPoints:          pointsToPlayerMap,
					CurrentPlayerPoints:   g.Desk.GetTotalPoints(),
					RoundFinished:         g.RoundFinished,
					PlayerIDsThatAccepted: playerIDsThatAccepted,
					Players:               g.Players,
					User: model.User{
						ID:   playerID,
						Name: player.Name,
					},
				})
			} else {
				// notify other players that the player reconnected
				server.server.Send(impl.StructMessageResponse(responses.PlayerConnectionChanged{
					PlayerID:     playerID,
					Disconnected: false,
				}), player.ID, 0)
			}
		}
		return resp
	} else {
		state, _ := server.server.Router.UserStates[playerID]
		log.Errorf("Player ID=%d reconnected but his game could not be found. Player state was %d", playerID, state.Id())
		server.server.Router.UserStates[playerID] = AuthorizedState{}
		return nil
	}
}

// on player willingly leaving the game
func (server *GameServer) OnPlayerLeavingGame(playerID int) def.Response {
	log.Infof("OnPlayerLeavingGame event by player ID=%d", playerID)
	state, exists := server.server.Router.UserStates[playerID]
	if exists {
		server.PlayerLeft(playerID, state.Id(), true)
	}

	delete(server.gamesByPlayerID, playerID)
	server.server.Router.IgnoreTransitionStateChange = false
	return impl.SuccessResponse("Successfully left the game")
}

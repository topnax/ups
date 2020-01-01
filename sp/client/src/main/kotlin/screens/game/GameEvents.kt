package screens.game

import model.game.Letter
import model.game.Tile
import networking.messages.GameStartedResponse
import networking.messages.GameStateRegenerationResponse
import tornadofx.EventBus
import tornadofx.FXEvent

// game events

class NewLetterSackEvent(val letters: List<Letter>) : FXEvent(EventBus.RunOn.BackgroundThread)

class DeskChange(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class TileWithLetterClicked(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class LetterPlacedEvent(val letter: Letter, val letterView: LetterView) : FXEvent(EventBus.RunOn.BackgroundThread)

class TileSelectedEvent(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class PlayerStateChangedEvent() : FXEvent(EventBus.RunOn.BackgroundThread)

class GameStartedEvent(val message: GameStartedResponse) : FXEvent(EventBus.RunOn.BackgroundThread)

class RoundFinishedEvent() : FXEvent(EventBus.RunOn.BackgroundThread)

class GameStateRegenerationEvent(val response: GameStateRegenerationResponse) : FXEvent(EventBus.RunOn.BackgroundThread)

class ConnectionStateChanged(val connected: Boolean) : FXEvent(EventBus.RunOn.BackgroundThread)

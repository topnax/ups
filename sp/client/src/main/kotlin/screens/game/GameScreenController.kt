package screens.game

import javafx.application.Platform
import javafx.scene.control.Alert
import model.game.Desk
import model.game.Letter
import model.game.Tile
import model.lobby.Player
import mu.KotlinLogging
import networking.Network
import networking.messages.*
import tornadofx.*

private val logger = KotlinLogging.logger { }

class GameScreenController : Controller() {

    var selectedTile: Tile? = null

    val desk = Desk(Desk.getTilesFromTileTypes())

    var letters: MutableList<Letter> = mutableListOf()

    val placedLetters: MutableList<Letter> = mutableListOf()

    lateinit var gameView: GameView

    var players: List<Player> = emptyList()

    var activePlayerID = -1

    fun init(gameView: GameView) {
        this.gameView = gameView
    }

    fun onDock() {
        Network.getInstance().addMessageListener(::onTileUpdated)
        Network.getInstance().addMessageListener(::onTilesUpdated)
    }

    fun onUndock() {
        Network.getInstance().removeMessageListener(::onTileUpdated)
        Network.getInstance().removeMessageListener(::onTilesUpdated)

    }

    fun onTileUpdated(response: TileUpdatedResponse) {
        desk.tiles[response.tile.row][response.tile.column] = response.tile
        if (!response.tile.set) {
            desk.tiles[response.tile.row][response.tile.column].letter = null
        }
        fire(DeskChange(desk.tiles[response.tile.row][response.tile.column]))
    }

    fun onTilesUpdated(response: TilesUpdatedResponse) {
        for (tile in response.tiles) {
            desk.tiles[tile.row][tile.column] = tile
            if (!tile.set) {
                desk.tiles[tile.row][tile.column].letter = null
            }
            fire(DeskChange(desk.tiles[tile.row][tile.column]))
        }
    }

    init {

        subscribe<GameStartedEvent> {
            fire(NewLetterSackEvent(it.message.letters))
            activePlayerID = it.message.activePlayerId
            players = it.message.players
            fire(PlayerStateChangedEvent())
        }

        subscribe<TileSelectedEvent> { event ->
            if (activePlayerID != Network.User.id) {
                alert(Alert.AlertType.ERROR, "Nejste na tahu")
            } else {
                logger.debug { "Tile ${event.tile} selected" }
                selectedTile?.let {
                    it.selected = false
                    fire(DeskChange(it))
                }
                selectedTile = event.tile
                selectedTile?.selected = true
                fire(DeskChange(event.tile))
            }
        }

        subscribe<LetterPlacedEvent> {
            if (activePlayerID != Network.User.id) {
                Platform.runLater {
                    alert(Alert.AlertType.ERROR, "Nejste na tahu")
                }
            } else {
                Network.getInstance().send(LetterPlacedMessage(it.letter, selectedTile!!.column, selectedTile!!.row), { am ->
                    run {
                        if (am is SuccessResponseMessage) {
                            Platform.runLater {
                                it.letterView.removeFromParent()
                                logger.debug { "Letter of value ${it.letter} has been pressed GO!" }
                                val tile = desk.tiles[selectedTile!!.row][selectedTile!!.column]
                                selectedTile = null
                                tile.letter = it.letter
                                tile.letter?.run {
                                    letters.remove(it.letter)
                                    placedLetters.add(it.letter)
                                    tile.selected = false
                                    fire(DeskChange(tile))
                                }
                            }
                        }
                    }
                }
                )
            }
        }

        subscribe<NewLetterSackEvent> {
            letters = it.letters.toMutableList()
        }

        subscribe<TileWithLetterClicked> { event ->
            if (activePlayerID == Network.User.id) {
                if (placedLetters.contains(event.tile.letter)) {
                    Network.getInstance().send(LetterRemovedMessage(event.tile.column, event.tile.row), {
                        run {
                            placedLetters.remove(event.tile.letter)
                            letters.add(event.tile.letter!!)
                            desk.tiles[event.tile.row][event.tile.column].letter = null
                            fire(DeskChange(event.tile))
                            fire(NewLetterSackEvent(letters))
                        }
                    })
                }
            } else {
                alert(Alert.AlertType.ERROR, "Nejste na tahu")
            }
        }
    }
}

class NewLetterSackEvent(val letters: List<Letter>) : FXEvent(EventBus.RunOn.BackgroundThread)

class DeskChange(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class TileWithLetterClicked(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class LetterPlacedEvent(val letter: Letter, val letterView: LetterView) : FXEvent(EventBus.RunOn.BackgroundThread)

class TileSelectedEvent(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class PlayerStateChangedEvent() : FXEvent(EventBus.RunOn.BackgroundThread)

class GameStartedEvent(val message: GameStartedResponse) : FXEvent(EventBus.RunOn.BackgroundThread)


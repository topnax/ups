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

    var roundFinished = false

    val playerIdsWhoAcceptedWords = mutableListOf<Int>()

    fun init(gameView: GameView) {
        this.gameView = gameView
    }

    fun onDock() {
        Network.getInstance().addMessageListener(::onTileUpdated)
        Network.getInstance().addMessageListener(::onTilesUpdated)
        Network.getInstance().addMessageListener(::onNewRound)
        Network.getInstance().addMessageListener(::onRoundFinished)
        Network.getInstance().addMessageListener(::onPlayerAcceptedRound)
        Network.getInstance().addMessageListener(::onNewRoundResponse)
        Network.getInstance().addMessageListener(::onYourNewRoundResponse)
    }

    fun onUndock() {
        Network.getInstance().removeMessageListener(::onTileUpdated)
        Network.getInstance().removeMessageListener(::onTilesUpdated)
        Network.getInstance().removeMessageListener(::onNewRound)
        Network.getInstance().removeMessageListener(::onRoundFinished)
        Network.getInstance().removeMessageListener(::onPlayerAcceptedRound)
        Network.getInstance().removeMessageListener(::onYourNewRoundResponse)
    }

    fun onYourNewRoundResponse(response: YourNewRoundResponse) {
        activePlayerID = Network.User.id
        fire(PlayerStateChangedEvent())
        fire(NewLetterSackEvent(response.letters))
    }

    fun onNewRoundResponse(response: NewRoundResponse) {
        activePlayerID = response.activePlayerId
        fire(PlayerStateChangedEvent())
    }

    fun onPlayerAcceptedRound(response: PlayerAcceptedRoundResponse) {
        playerIdsWhoAcceptedWords.add(response.playerId)
        fire(PlayerStateChangedEvent())
    }

    fun onRoundFinished(response: RoundFinishedResponse) {
        roundFinished = true
        fire(RoundFinishedEvent())
    }

    fun onNewRound(response: NewRoundResponse) {
        roundFinished = false
        activePlayerID = response.activePlayerId
        fire(PlayerStateChangedEvent())
        gameView.finishButton.visibleProperty().set(activePlayerID == Network.User.id)
    }

    fun onTileUpdated(response: TileUpdatedResponse) {
        desk.tiles[response.tile.row][response.tile.column] = response.tile
        if (!response.tile.set) {
            desk.tiles[response.tile.row][response.tile.column].letter = null
        }
        fire(DeskChange(desk.tiles[response.tile.row][response.tile.column]))
    }

    fun onTilesUpdated(response: TilesUpdatedResponse) {
        for (tileRow in desk.tiles) {
            for (tile in tileRow) {
                tile.highlighted = false
                fire(DeskChange(tile))
            }
        }

        for (tile in response.tiles) {
            desk.tiles[tile.row][tile.column] = tile
            if (!tile.set) {
                desk.tiles[tile.row][tile.column].letter = null
            }
            fire(DeskChange(desk.tiles[tile.row][tile.column]))
        }
    }

    fun onFinishRoundButtonClicked() {
        if (activePlayerID == Network.User.id) {
            Network.getInstance().send(FinishRoundMessage(), {
                Platform.runLater {
                    gameView.finishButton.visibleProperty().set(false)
                    alert(Alert.AlertType.INFORMATION, "Waiting for confirmation from other players :)")
                }
            })
        }
    }

    fun onAcceptWordsButtonClicked() {
        logger.debug { "Accept words button clicked" }
        if (activePlayerID != Network.User.id) {
            logger.debug { "The player accepting words is not the active player" }
            if (roundFinished) {
                logger.debug { "Round finished, sending approval..." }
                Network.getInstance().send(ApproveWordsMessage(), {
                    Platform.runLater {
                        alert(Alert.AlertType.INFORMATION, "Successfully approved :)")
                    }
                })
            }
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
                Platform.runLater {
                    alert(Alert.AlertType.ERROR, "Nejste na tahu")
                }
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

class RoundFinishedEvent() : FXEvent(EventBus.RunOn.BackgroundThread)


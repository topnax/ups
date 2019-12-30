package screens.game

import javafx.application.Platform
import javafx.scene.control.Alert
import model.game.Desk
import model.game.Letter
import model.game.Tile
import model.lobby.Player
import mu.KotlinLogging
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import screens.DisconnectedEvent
import screens.ServerRestartedEvent
import screens.ServerRestartedUnauthorizedEvent
import screens.disconnected.DisconnectedScreenView
import screens.initial.InitialScreenView
import screens.mainmenu.MainMenuView
import tornadofx.*

private val logger = KotlinLogging.logger { }

class GameScreenController : Controller(), ConnectionStatusListener {

    val previouslyUpdatedTiles = mutableListOf<Tile>()

    var selectedTile: Tile? = null

    val desk = Desk(Desk.getTilesFromTileTypes())

    var letters: MutableList<Letter> = mutableListOf()

    val placedLetters: MutableList<Letter> = mutableListOf()

    lateinit var gameView: GameView

    var players: List<Player> = emptyList()

    var activePlayerID = -1

    var currentRoundPlayerPoints = 0

    var roundFinished = false

    var wordsAccepted = false

    val playerIdsWhoAcceptedWords = mutableListOf<Int>()

    val playerPointsMap = mutableMapOf<Int, Int>()

    fun  init(gameView: GameView) {
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
        Network.getInstance().addMessageListener(::onPlayerDeclinedWordsResponse)
        Network.getInstance().addMessageListener(::onGameEndedResponse)
        Network.getInstance().addMessageListener(::onPlayerConnectionChangedResponse)
        Network.getInstance().addMessageListener(::onGameStateRegenerationResponse)
        reset()
    }

    fun onUndock() {
        Network.getInstance().removeMessageListener(::onTileUpdated)
        Network.getInstance().removeMessageListener(::onTilesUpdated)
        Network.getInstance().removeMessageListener(::onNewRound)
        Network.getInstance().removeMessageListener(::onRoundFinished)
        Network.getInstance().removeMessageListener(::onPlayerAcceptedRound)
        Network.getInstance().removeMessageListener(::onNewRoundResponse)
        Network.getInstance().removeMessageListener(::onYourNewRoundResponse)
        Network.getInstance().removeMessageListener(::onPlayerDeclinedWordsResponse)
        Network.getInstance().removeMessageListener(::onGameEndedResponse)
        Network.getInstance().removeMessageListener(::onPlayerConnectionChangedResponse)
        Network.getInstance().removeMessageListener(::onGameStateRegenerationResponse)
    }

    private fun reset() {
        previouslyUpdatedTiles.clear()
        selectedTile = null
        placedLetters.clear()
        currentRoundPlayerPoints = 0
        roundFinished = false
        wordsAccepted = false
        playerIdsWhoAcceptedWords.clear()
        playerPointsMap.clear()
    }

    fun onGameStateRegenerationResponse(response: GameStateRegenerationResponse) {
        fire(GameStateRegenerationEvent(response))
    }

    fun onGameEndedResponse(response: GameEndedResponse) {
        val maxPoints = response.playerPoints.map { it.key.toInt() }.max()
        val text = response.playerPoints.map {
            "${it.value.name} - ${it.key} pts " + if (it.key.toInt() == maxPoints) "(WINNER)" else ""
        }.joinToString(separator = "\n")
        Platform.runLater {
            alert(Alert.AlertType.INFORMATION, "Game has ended", text)
            gameView.replaceWith<MainMenuView>()
        }
    }

    fun onPlayerConnectionChangedResponse(response: PlayerConnectionChangedResponse) {
        for (player in players) {
            if (player.id == response.playerId) {
                player.disconnected = response.disconnected
                fire(PlayerStateChangedEvent())
                break
            }
        }
    }

    fun onYourNewRoundResponse(response: YourNewRoundResponse) {
        wordsAccepted = false
        roundFinished = false
        currentRoundPlayerPoints = 0

        activePlayerID = Network.User.id
        playerIdsWhoAcceptedWords.clear()
        fire(PlayerStateChangedEvent())
        fire(NewLetterSackEvent(response.letters))
        for (tile in previouslyUpdatedTiles) {
            tile.highlighted = false
            fire(DeskChange(tile))
        }
        previouslyUpdatedTiles.clear()
    }

    fun onNewRoundResponse(response: NewRoundResponse) {
        currentRoundPlayerPoints = 0
        wordsAccepted = false
        roundFinished = false
        playerIdsWhoAcceptedWords.clear()
        activePlayerID = response.activePlayerId
        fire(PlayerStateChangedEvent())
    }

    fun onPlayerDeclinedWordsResponse(response: PlayerDeclinedWordsResponse) {
        wordsAccepted = false
        roundFinished = false
        playerIdsWhoAcceptedWords.clear()
        fire(PlayerStateChangedEvent())
        Platform.runLater {
            alert(Alert.AlertType.ERROR, "${response.playerName} #${response.playerId} has declined the words")
        }
    }

    fun onPlayerAcceptedRound(response: PlayerAcceptedRoundResponse) {
        playerIdsWhoAcceptedWords.add(response.playerId)
        fire(PlayerStateChangedEvent())
    }

    fun onRoundFinished(response: RoundFinishedResponse) {
        roundFinished = true
        wordsAccepted = false
        fire(RoundFinishedEvent())
    }

    fun onNewRound(response: NewRoundResponse) {
        currentRoundPlayerPoints = 0
        roundFinished = false
        wordsAccepted = false
        playerIdsWhoAcceptedWords.clear()
        activePlayerID = response.activePlayerId
        fire(PlayerStateChangedEvent())
        gameView.finishButton.visibleProperty().set(activePlayerID == Network.User.id)

        for (tile in previouslyUpdatedTiles) {
            tile.highlighted = false
            fire(DeskChange(tile))
        }
        previouslyUpdatedTiles.clear()
    }

    fun onTileUpdated(response: TileUpdatedResponse) {
        desk.tiles[response.tile.row][response.tile.column] = response.tile
        if (!response.tile.set) {
            desk.tiles[response.tile.row][response.tile.column].letter = null
        }
        fire(DeskChange(desk.tiles[response.tile.row][response.tile.column]))
    }

    fun onTilesUpdated(response: TilesUpdatedResponse) {
        for (tile in previouslyUpdatedTiles) {
            tile.highlighted = false
        }

        for (tile in response.tiles) {
            desk.tiles[tile.row][tile.column] = tile
            if (!tile.set) {
                desk.tiles[tile.row][tile.column].letter = null
            }
            fire(DeskChange(desk.tiles[tile.row][tile.column]))
        }

        for (tile in previouslyUpdatedTiles) {
            // TODO might remember which tiles were updated through the response
            fire(DeskChange(desk.tiles[tile.row][tile.column]))
        }
        previouslyUpdatedTiles.clear()

        previouslyUpdatedTiles.addAll(response.tiles)

        currentRoundPlayerPoints = response.currentPlayerPoints
        playerPointsMap[activePlayerID] = response.currentPlayerTotalPoints
        fire(PlayerStateChangedEvent())
    }

    fun onFinishRoundButtonClicked() {
        logger.debug { "onFinishRoundButtonClicked() => activePlayerID=$activePlayerID, Network.User.id=${Network.User.id}" }
        if (activePlayerID == Network.User.id) {
            Network.getInstance().send(FinishRoundMessage(), {
                if (it !is NewRoundResponse) {
                    roundFinished = true
                    Platform.runLater {
                        gameView.finishButton.visibleProperty().set(false)
                    }
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
                    if (it !is AcceptResultedInNewRound) {
                        wordsAccepted = true
                        playerIdsWhoAcceptedWords.add(Network.User.id)
                        fire(PlayerStateChangedEvent())
                    }
                })
            }
        }
    }

    fun onDeclineWordsClicked() {
        logger.debug { "Decline words button clicked" }
        if (activePlayerID != Network.User.id && roundFinished) {
            Network.getInstance().send(DeclineWordsMessage(), {
                wordsAccepted = false
                roundFinished = false
                playerIdsWhoAcceptedWords.clear()
                fire(PlayerStateChangedEvent())
            })
        }
    }

    init {
        players.forEach { playerPointsMap[it.id] = 0 }

        subscribe<DisconnectedEvent> {
            Platform.runLater {
                gameView.replaceWith<DisconnectedScreenView>()
            }
        }

        subscribe<GameStartedEvent> {
            fire(NewLetterSackEvent(it.message.letters))
            activePlayerID = it.message.activePlayerId
            players = it.message.players
            wordsAccepted = false
            roundFinished = false
            selectedTile = null
            playerIdsWhoAcceptedWords.clear()

            logger.info { "activePlayerID=$activePlayerID, networkUserId=${Network.User.id}" }
            players.forEach { playerPointsMap[it.id] = 0 }
            fire(PlayerStateChangedEvent())
        }

        subscribe<GameStateRegenerationEvent> {
            logger.warn { "Inside GameScreenController state regen" }
            activePlayerID = it.response.activePlayerId
            currentRoundPlayerPoints = it.response.currentPlayerPoints
            playerIdsWhoAcceptedWords.addAll(it.response.playerIdsThatAccepted)
            if (playerIdsWhoAcceptedWords.contains(Network.User.id)) {
                wordsAccepted = true
            }
            players = it.response.players
            players.forEach { playerPointsMap[it.id] = 0 }
            playerPointsMap.putAll(it.response.playerPoints.map { it.value.id to it.key.toInt() })
            roundFinished = it.response.roundFinished

            playerPointsMap[activePlayerID]?.let {
                playerPointsMap[activePlayerID] = it + currentRoundPlayerPoints
            }

            for (tile in it.response.tiles) {
                desk.tiles[tile.row][tile.column] = tile
                if (!tile.set) {
                    desk.tiles[tile.row][tile.column].letter = null
                } else {
                    logger.warn { "Setting a letter at ${tile.row} and c ${tile.column}" }
                }
                fire(DeskChange(desk.tiles[tile.row][tile.column]))
            }
            previouslyUpdatedTiles.addAll(it.response.tiles)
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

        subscribe<LetterPlacedEvent> { event ->
            if (activePlayerID != Network.User.id) {
                Platform.runLater {
                    alert(Alert.AlertType.ERROR, "Nejste na tahu")
                }
            } else {
                selectedTile?.let {
                    Network.getInstance().send(LetterPlacedMessage(event.letter, it.column, selectedTile!!.row), { am ->
                        run {
                            if (am is SuccessResponseMessage) {
                                Platform.runLater {
                                    event.letterView.removeFromParent()
                                    logger.debug { "Letter of value ${it.letter} has been pressed GO!" }
                                    val tile = desk.tiles[selectedTile!!.row][selectedTile!!.column]
                                    selectedTile = null
                                    tile.letter = event.letter
                                    tile.letter?.run {
                                        letters.remove(event.letter)
                                        placedLetters.add(event.letter)
                                        tile.selected = false
                                        fire(DeskChange(tile))
                                    }
                                }
                            }
                        }
                    })
                } ?: run {
                    Platform.runLater {
                        alert(Alert.AlertType.ERROR, "No tile selected")
                    }
                }
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
                            if (it is SuccessResponseMessage) {
                                selectedTile = event.tile
                                placedLetters.remove(event.tile.letter)
                                letters.add(event.tile.letter!!)
                                desk.tiles[event.tile.row][event.tile.column].letter = null
                                desk.tiles[event.tile.row][event.tile.column].selected = true
                                fire(NewLetterSackEvent(letters))
                                fire(DeskChange(desk.tiles[event.tile.row][event.tile.column]))
                            }
                        }
                    })
                }
            } else {
                alert(Alert.AlertType.ERROR, "It is not your turn")
            }
        }

        Network.getInstance().connectionStatusListeners.add(this)
    }

    override fun onConnected() {
        fire(ConnectionStateChanged(true))
    }

    override fun onUnreachable() {
        fire(ConnectionStateChanged(false))
    }

    override fun onFailedAttempt(attempt: Int) {
        fire(ConnectionStateChanged(false))
    }

    override fun onReconnected() {
        fire(ConnectionStateChanged(true))
    }

    private fun resetDesk() {
        for (row in desk.tiles) {
            for (tile in row) {
                tile.highlighted = false
                tile.letter = null
                tile.set = false
                fire(DeskChange(tile))
            }
        }
    }

    fun leaveGame() {
        Network.getInstance().send(LeaveGame(), {
            if (it is SuccessResponseMessage) {
                Platform.runLater {
                    gameView.replaceWith<MainMenuView>()
                }
            }
        })
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

class GameStateRegenerationEvent(val response: GameStateRegenerationResponse) : FXEvent(EventBus.RunOn.BackgroundThread)

class ConnectionStateChanged(val connected: Boolean) : FXEvent(EventBus.RunOn.BackgroundThread)


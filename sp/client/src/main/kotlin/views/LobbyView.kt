package views

import MainMenuView
import javafx.application.Platform
import javafx.geometry.Pos
import javafx.scene.control.Alert
import javafx.scene.control.Button
import javafx.scene.control.ListView
import model.lobby.Lobby
import model.lobby.Player
import networking.Network
import networking.messages.*
import tornadofx.*
import java.util.*
import kotlin.concurrent.schedule

class LobbyView : View() {

    private lateinit var playerListView: ListView<String>
    private lateinit var readyButton: Button
    private lateinit var startButton: Button

    var ready = false

    val lobby: Lobby by param(Lobby(listOf(), -1, Player("", -1, false)))
    val player: Player by param(Player("", -1, false))

    private fun onLobbyUpdated(message: LobbyUpdatedResponse) {
        Platform.runLater {
            update(message.lobby)
        }
    }

    private fun onLobbyDestroyed(message: LobbyDestroyedResponse) {
        Platform.runLater {
            val mainMenu = find<MainMenuView>()
            replaceWith(mainMenu)
            alert(Alert.AlertType.INFORMATION, "Lobby disbandned", "The lobby has been disbanded")
        }
    }

    override val root = vbox(spacing = 10) {
        alignment = Pos.CENTER
        padding = insets(10)
        text("Welcome to game lobby") {
            alignment = Pos.CENTER
        }
        playerListView = listview {}
        readyButton = button("Ready")
        readyButton.action {
            onReadyButtonClicked()
        }
        button("Leave lobby") {
            alignment = Pos.CENTER
            action {
                leaveLobby()
            }
        }
        startButton = button("Start button") {
            alignment = Pos.CENTER
            action {
                startLobby()
            }
            visibleProperty().set(false)
        }
    }

    private fun startLobby() {
        Network.getInstance().send(StartLobbyMessage(), { applicationMessage ->
            Platform.runLater {
                if (applicationMessage is LobbyStartedResponse) {
                    alert(Alert.AlertType.INFORMATION, "Lobby started")
                }
            }
        })
    }

    private fun onReadyButtonClicked() {
        Network.getInstance().send(PlayerReadyToggleMessage(!ready), { am ->
            run {
                when (am) {
                    is LobbyUpdatedResponse -> {
                        ready = !ready
                        Platform.runLater { update(am.lobby) }
                    }
                    else -> Platform.runLater { alert(Alert.AlertType.ERROR, "Error", "Could not toggle state") }
                }
            }
        })
    }

    private fun leaveLobby() {
        Network.getInstance().send(LeaveLobbyMessage(), ignoreErrors = true)
        val mainMenuView = find<MainMenuView>()
        replaceWith(mainMenuView)
        Timer().schedule(2000) {
            Platform.runLater {
                mainMenuView.controller.refreshLobbies()
            }
        }
    }

    fun update(lobby: Lobby) {
        playerListView.items.clear()
        lobby.players.forEach {
            var displayedName = if (lobby.owner.id == it.id) it.name + " (owner)" else it.name
            if (!it.ready) {
                displayedName += " not ready"
            }
            playerListView.items.add(displayedName)
        }
        if (player == lobby.owner) {
            startButton.visibleProperty().set(true)
        } else {
            startButton.visibleProperty().set(false)
        }

        startButton.disableProperty().set(lobby.players.filter { it.ready }.count() != lobby.players.size || lobby.players.size < 2)

    }

    override fun onDock() {
        ready = false
        update(lobby)
        Network.getInstance().addMessageListener(::onLobbyUpdated)
        Network.getInstance().addMessageListener(::onLobbyDestroyed)
        if (Network.User.id != lobby.owner.id) {
            Network.getInstance().addMessageListener(::onLobbyStarted)
        }
    }

    private fun onLobbyStarted(lobbyStartedResponse: LobbyStartedResponse) {
        Platform.runLater {
            alert(Alert.AlertType.INFORMATION, "Leader has started the lobby")
        }
    }

    override fun onUndock() {
        super.onUndock()
        Network.getInstance().removeMessageListener(::onLobbyDestroyed)
        Network.getInstance().removeMessageListener(::onLobbyUpdated)
        if (Network.User.id != lobby.owner.id) {
            Network.getInstance().removeMessageListener(::onLobbyStarted)
        }
    }
}
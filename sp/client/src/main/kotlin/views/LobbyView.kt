package views

import MainMenuView
import javafx.application.Platform
import javafx.geometry.Pos
import javafx.scene.control.Alert
import javafx.scene.control.ListView
import model.lobby.Lobby
import model.lobby.Player
import networking.Network
import networking.messages.LeaveLobbyMessage
import networking.messages.LobbyDestroyedResponse
import networking.messages.LobbyJoinedMessage
import networking.messages.PlayerJoinedLobby
import tornadofx.*
import java.util.*
import kotlin.concurrent.schedule

class LobbyView : View() {

    private lateinit var playerListView: ListView<String>

    val lobby: Lobby by param(Lobby(listOf(), -1, Player("", -1)))

    private fun onLobbyUpdated(message: LobbyJoinedMessage) {
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
        button("Leave lobby") {
            alignment = Pos.CENTER
            action {
                leaveLobby()
            }
        }
    }

    private fun leaveLobby() {
        Network.getInstance().send(LeaveLobbyMessage())
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
            playerListView.items.add(if (lobby.owner.id == it.id) it.name + " (owner)" else it.name)
        }
    }

    override fun onDock() {
        update(lobby)
        Network.getInstance().addMessageListener(::onLobbyUpdated)
        Network.getInstance().addMessageListener(::onLobbyDestroyed)
    }

    override fun onUndock() {
        super.onUndock()
        Network.getInstance().removeMessageListener(::onLobbyDestroyed)
        Network.getInstance().removeMessageListener(::onLobbyUpdated)
    }
}
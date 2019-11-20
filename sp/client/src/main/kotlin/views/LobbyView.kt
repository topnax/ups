package views

import javafx.application.Platform
import javafx.scene.control.ListView
import networking.Network
import networking.messages.PlayerJoinedLobby
import tornadofx.View
import tornadofx.listview
import tornadofx.text
import tornadofx.vbox

class LobbyView : View() {

    lateinit var playerListView: ListView<String>

    val playerName: String by param("Unknown name")

    init {
        Network.getInstance().addMessageListener(::onPlayerJoinedLobby)
    }

    private fun onPlayerJoinedLobby(message: PlayerJoinedLobby) {
        Platform.runLater {
            println("Player of name " + message.playerName + " has joined")
            playerListView.items.add(message.playerName)
            Network.getInstance().removeMessageListener(::onPlayerJoinedLobby)
        }
    }

    override val root = vbox {
        text("Welcome to game lobby")
        playerListView = listview {
            items.add(playerName)
        }
    }
}
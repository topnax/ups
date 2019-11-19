package views

import javafx.application.Platform
import javafx.scene.control.ListView
import networking.Network
import networking.messages.ApplicationMessage
import networking.messages.ErrorResponseMessage
import networking.messages.PlayerJoinedLobby
import networking.messages.SuccessResponseMessage
import tornadofx.View
import tornadofx.listview
import tornadofx.text
import tornadofx.vbox

class LobbyView : View() {

    lateinit var playerListView: ListView<String>

    val playerName: String by param("Unknown name")

    init {
        Network.getInstance().addMessageListener(PlayerJoinedLobby::class.java) { am: ApplicationMessage->
            run {
                if (am is PlayerJoinedLobby) {
                    Platform.runLater{
                        println("Player of name " + am.playerName +" has joined")
                        playerListView.items.add(am.playerName)
                    }
                }
            }
        }
    }

    override val root = vbox {
        text("Welcome to game lobby")
        playerListView = listview {
            items.add(playerName)
        }
    }
}
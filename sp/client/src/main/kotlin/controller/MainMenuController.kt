package controller

import MainMenuView
import javafx.application.Platform
import javafx.collections.ObservableList
import model.lobby.Lobby
import model.lobby.LobbyViewModel
import model.lobby.Player
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import networking.reader.MessageReader
import networking.receiver.Message
import tornadofx.Controller
import tornadofx.observableList
import views.LobbyView

class MainMenuController : Controller(), ConnectionStatusListener {

    lateinit var mainMenuView: MainMenuView

    var lobbyViewModels: ObservableList<LobbyViewModel> = observableList()

    fun init(mainMenuView: MainMenuView) {
        this.mainMenuView = mainMenuView
        Network.getInstance().connectionStatusListeners.add(this)
        connectTo("localhost", 10000)
    }

    override fun onConnected() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
            mainMenuView.serverMenu.text = "Connected to ${Network.getInstance().tcpLayer?.hostname}"
        }
        refreshLobbies()
    }

    override fun onUnreachable() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
            mainMenuView.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} is unreachable"
        }
    }

    override fun onFailedAttempt(attempt: Int) {
        Platform.runLater {
            mainMenuView.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} did not respond. Attempt $attempt"
        }
    }

    fun newLobby(name: String) {
        Network.getInstance().send(CreateLobbyMessage(name), { am: ApplicationMessage ->
            run {
                when (am) {
                    is SuccessResponseMessage -> {
                        Platform.runLater {
                            val player = Player(name, -1)
                            mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::lobby to Lobby(listOf(player), -1, player))))
                        }
                    }
                    is ErrorResponseMessage -> println("fail ${am.content}")
                }
            }
        })
    }

    fun refreshLobbies() {
        Network.getInstance().send(GetLobbiesMessage(), { am: ApplicationMessage ->
            run {
                when (am) {
                    is GetLobbiesResponse -> {
                        lobbyViewModels.clear()
                        am.lobbies.forEach {
                            lobbyViewModels.add(LobbyViewModel(it.id, it.owner.name, it.players.size))
                        }
                    }
                    else -> lobbyViewModels.clear()
                }
            }
        })
    }

    fun connectTo(hostname: String, port: Int) {
        mainMenuView.setNetworkElementsEnabled(false)
        Network.getInstance().connectTo(hostname, port)
    }

    fun onJoinLobby(id: Int) {
        Network.getInstance().send(JoinLobbyMessage(id, mainMenuView.nameTextField.text), { am: ApplicationMessage ->
            run {
                Platform.runLater {
                    if (am is LobbyJoinedMessage) {
                        println("Owner is of id #${am.lobby.owner.id}")
                        am.lobby.players.forEach {
                            println("player ${it.name} of id ${it.id}")
                        }
                        mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::lobby to am.lobby)))
                    }
                }
            }
        })
    }

}

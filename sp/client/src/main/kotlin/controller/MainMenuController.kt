package controller

import MainMenuView
import javafx.application.Platform
import javafx.collections.ObservableList
import javafx.scene.control.Alert
import model.lobby.LobbyViewModel
import mu.KotlinLogging
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import tornadofx.Controller
import tornadofx.alert
import tornadofx.observableList
import views.LobbyView

private val logger = KotlinLogging.logger { }

class MainMenuController : Controller(), ConnectionStatusListener {

    lateinit var mainMenuView: MainMenuView

    var lobbyViewModels: ObservableList<LobbyViewModel> = observableList()

    fun init(mainMenuView: MainMenuView) {
        this.mainMenuView = mainMenuView

        Network.getInstance().connectionStatusListeners.add(this)

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

    fun newLobby() {
        Network.getInstance().send(CreateLobbyMessage(), { am: ApplicationMessage ->
            run {
                when (am) {
                    is LobbyJoinedResponse -> {
                        Platform.runLater {
                            mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::lobby to am.lobby, LobbyView::player to am.player)))
                        }
                    }
                    is ErrorResponseMessage -> logger.error { "Could not create a new lobby. Response content '${am.content}'" }
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

    fun onJoinLobby(id: Int) {
        Network.getInstance().send(JoinLobbyMessage(id), { am: ApplicationMessage ->
            run {
                Platform.runLater {
                    if (am is LobbyJoinedResponse) {
                        mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::lobby to am.lobby)))
                    }
                }
            }
        })
    }
}

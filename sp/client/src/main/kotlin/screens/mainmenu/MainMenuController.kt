package screens.mainmenu


import javafx.application.Platform
import javafx.collections.ObservableList
import model.lobby.Lobby
import model.lobby.LobbyViewModel
import mu.KotlinLogging
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import screens.DisconnectedEvent
import screens.disconnected.DisconnectedScreenView
import screens.lobby.LobbyView
import tornadofx.Controller
import tornadofx.observableList

private val logger = KotlinLogging.logger { }

class MainMenuController : Controller(), ConnectionStatusListener {

    private lateinit var mainMenuView: MainMenuView

    var lobbyViewModels: ObservableList<LobbyViewModel> = observableList()

    init {
        subscribe<DisconnectedEvent> {
            Platform.runLater {
                mainMenuView.replaceWith<DisconnectedScreenView>()
            }
        }
    }

    fun init(mainMenuView: MainMenuView) {
        this.mainMenuView = mainMenuView

        Network.getInstance().connectionStatusListeners.add(this)
    }

    override fun onReconnected() {}

    override fun onConnected() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
        }
        refreshLobbies()
    }

    override fun onUnreachable() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
        }
    }

    override fun onFailedAttempt(attempt: Int) {}

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

            Platform.runLater {
                when (am) {
                    is GetLobbiesResponse -> {
                        updateLobbiesTable(am.lobbies)
                    }
                    else -> lobbyViewModels.clear()
                }
            }

        }, ignoreErrors = true)
    }

    fun updateLobbiesTable(lobbies: List<Lobby>) {
        lobbyViewModels.clear()
        lobbies.forEach {
            lobbyViewModels.add(LobbyViewModel(it.id, it.owner.name, it.players.size))
        }
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

package controller

import MainMenuView
import javafx.application.Platform
import javafx.collections.ObservableList
import javafx.scene.control.Alert
import model.lobby.Lobby
import model.lobby.LobbyViewModel
import model.lobby.Player
import mu.KotlinLogging
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import networking.reader.MessageReader
import networking.receiver.Message
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
        mainMenuView.primaryStage.setOnCloseRequest {
            logger.debug { "Primary stage closing" }
            Network.getInstance().stop()
        }

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

    fun newLobby() {
        if (validateName()) {
            val name = mainMenuView.nameTextField.text.trim()
            Network.getInstance().send(CreateLobbyMessage(name), { am: ApplicationMessage ->
                run {
                    when (am) {
                        is SuccessResponseMessage -> {
                            Platform.runLater {
                                val player = Player(name, -1, false)
                                mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::lobby to Lobby(listOf(player), -1, player))))
                            }
                        }
                        is ErrorResponseMessage -> logger.error { "Could not create a new lobby. Response content '${am.content}'" }
                    }
                }
            })
        }
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
        if (validateName()) {
            Network.getInstance().send(JoinLobbyMessage(id, mainMenuView.nameTextField.text.trim()), { am: ApplicationMessage ->
                run {
                    Platform.runLater {
                        if (am is LobbyJoinedMessage) {
                            mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::lobby to am.lobby)))
                        }
                    }
                }
            })
        }
    }

    private fun validateName(): Boolean {
        if (mainMenuView.nameTextField.text.trim().isNotEmpty()) {
            return true
        }
        alert(Alert.AlertType.ERROR, "Error", "Name must be of non-zero length")
        return false
    }
}

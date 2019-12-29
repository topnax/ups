package screens.mainmenu

import javafx.application.Platform
import javafx.scene.control.Button
import javafx.scene.control.Label
import javafx.scene.control.MenuItem
import javafx.scene.layout.Priority
import model.lobby.LobbyViewModel
import networking.Network
import networking.messages.GetLobbiesResponse
import tornadofx.*
import java.util.*
import kotlin.concurrent.schedule

class MainMenuView : View() {

    fun setNetworkElementsEnabled(b: Boolean) {
        createLobbyButton.disableProperty().set(!b)
    }

    private lateinit var createLobbyButton: Button
    lateinit var serverMenu: MenuItem

    val controller: MainMenuController by inject()

    override val root = borderpane {

        center = vbox(spacing = 10.0) {
            padding = insets(10)
            prefWidth = 10.0
            hbox(spacing = 10.0) {
                createLobbyButton = button("Create a lobby")
                createLobbyButton.action {
                    controller.newLobby()
                }

                button("Refresh lobby list").action {
                    controller.refreshLobbies()
                }

                label (Network.User.name)
            }

            tableview(controller.lobbyViewModels) {
                placeholder = Label("No lobbies")
                insets(10.0)
                column("ID", LobbyViewModel::idProperty)
                column("Owner", LobbyViewModel::ownerName)
                column("Players", LobbyViewModel::playersProperty)
                vboxConstraints {
                    vGrow = Priority.ALWAYS
                }
                onDoubleClick {
                    this.selectedItem?.let {
                        controller.onJoinLobby(it.id)
                    }
                }

            }
        }
        controller.init(this@MainMenuView)
    }

    override fun onDock() {
        super.onDock()
        Timer().schedule(1000) {
            if (isDocked) {
                controller.refreshLobbies()
            }
        }

        Network.getInstance().addMessageListener(::onLobbyUpdated)
    }

    override fun onUndock() {
        super.onUndock()
        Network.getInstance().removeMessageListener(::onLobbyUpdated)
    }


    private fun onLobbyUpdated(message: GetLobbiesResponse) {
        Platform.runLater {
            controller.updateLobbiesTable(message.lobbies)
        }
    }

}



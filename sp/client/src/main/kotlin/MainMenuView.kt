import controller.MainMenuController
import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import javafx.scene.control.Button
import javafx.scene.control.Label
import javafx.scene.control.MenuItem
import javafx.scene.control.TextField
import javafx.scene.layout.Priority
import model.lobby.LobbyViewModel
import tornadofx.*

class MainMenuView : View() {

    fun setNetworkElementsEnabled(b: Boolean) {
        createLobbyButton.disableProperty().set(!b)
        serverMenu.disableProperty().set(!b)
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
        controller.refreshLobbies()
    }
}



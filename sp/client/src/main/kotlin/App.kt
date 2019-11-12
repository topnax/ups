import controller.MainMenuController
import javafx.scene.control.Label
import javafx.scene.control.MenuItem
import javafx.scene.layout.Priority
import model.lobby.Lobby
import tornadofx.*

class MainMenuView : View() {

    lateinit var serverMenu: MenuItem

    val controller: MainMenuController by inject()

    override val root = borderpane {

        top = menubar {
            serverMenu = menu("127.0.0.1") {
                item("Change server") {
                    action {
                        serverMenu.text = "127.1.1.1"
                        controller.newLobby()
                    }
                }
            }
            menu("Options") {
                item("Change foo")
                item("Change bar")
            }
        }

        center = vbox(spacing = 10.0) {
            padding = insets(10)
            prefWidth = 10.0
            hbox(spacing = 10.0) {
                button("Create a lobby").action {
                    serverMenu.text = "127.1.1.9"
                    controller.newLobby()
                }

                button("Join").action {
                    replaceWith<GameView>()
                }

            }

            tableview(controller.lobbies) {
                placeholder = Label("No lobbies found")
                insets(10.0)
                column("ID", Lobby::idProperty)
                column("Owner", Lobby::owner)
                column("Players", Lobby::playersProperty)
                vboxConstraints {
                    vGrow = Priority.ALWAYS
                }
            }
        }
    }
}

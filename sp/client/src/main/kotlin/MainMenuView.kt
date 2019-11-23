import controller.MainMenuController
import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import javafx.scene.control.Button
import javafx.scene.control.Label
import javafx.scene.control.MenuItem
import javafx.scene.control.TextField
import javafx.scene.layout.Priority
import model.lobby.LobbyViewModel
import networking.Network
import networking.messages.ApplicationMessage
import networking.messages.GetLobbiesMessage
import networking.messages.GetLobbiesResponse
import org.slf4j.LoggerFactory

import tornadofx.*

class MainMenuView : View() {


    fun setNetworkElementsEnabled(b: Boolean) {
        createLobbyButton.disableProperty().set(!b)
        serverMenu.disableProperty().set(!b)
    }

    private lateinit var createLobbyButton: Button
    lateinit var nameTextField: TextField
    lateinit var serverMenu: MenuItem

    val controller: MainMenuController by inject()

    override val root = borderpane {
        top = menubar {
            serverMenu = menu("127.0.0.1") {
                item("Change server") {
                    action {
                        val modal = find<JoinServerView>()
                        modal.mainViewController = controller
                        modal.openModal()
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
                createLobbyButton = button("Create a lobby")
                createLobbyButton.action {
                    controller.newLobby()
                }

                nameTextField = textfield {}

                button("Join TEST").action {
                    controller.refreshLobbies()
                }

            }

            tableview(controller.lobbyViewModels) {
                placeholder = Label("No lobbyViewModels found")
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

    class JoinServerView: View() {
        var mainViewController: MainMenuController? = null
        val model = ConnectMetaModel(ConnectMeta())
        override val root = vbox(spacing = 10) {
            form {
                fieldset("New connection") {
                    field("Hostname") {
                        textfield(model.hostName).validator {
                            if (it.isNullOrBlank()) error("Not a valid hostname") else null
                        }
                    }
                    field("Port") {
                        textfield(model.port).filterInput {
                            it.controlNewText.isInt() && it.controlNewText.toInt() > 0 && it.controlNewText.toInt() < 65535
                        }
                    }
                    hbox(spacing = 10) {
                        button("Save") {
                            enableWhen(model.valid)
                            action {
                                save()
                            }
                        }
                        button("Reset").action {
                            model.rollback()
                        }
                    }

                }
            }

        }

        private fun save() {
            model.commit()

            mainViewController?.connectTo(model.hostName.value, model.port.value.toInt())
            close()
        }

        inner class ConnectMeta(hostName: String? = null, port: Int? = null) {
            val hostNameProperty = SimpleStringProperty(this, "hostName", hostName)
            var name by hostNameProperty

            val portProperty = SimpleIntegerProperty(this, "port", 10000)
            var port by portProperty
        }

        inner class ConnectMetaModel(connectMeta: ConnectMeta) : ItemViewModel<ConnectMeta>(connectMeta) {
            val hostName = bind(ConnectMeta::hostNameProperty)
            val port = bind(ConnectMeta::portProperty)
        }
    }
}



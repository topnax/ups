import controller.MainMenuController
import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import javafx.scene.Parent
import javafx.scene.control.Button
import javafx.scene.control.Label
import javafx.scene.control.MenuItem
import javafx.scene.layout.Priority
import model.lobby.Lobby
import networking.Network
import networking.messages.GetLobbiesMessage
import tornadofx.*
import java.util.regex.Pattern

class MainMenuView : View() {


    fun setNetworkElementsEnabled(b: Boolean) {
        createLobbyButton.disableProperty().set(!b)
        serverMenu.disableProperty().set(!b)
    }

    private lateinit var createLobbyButton: Button
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
                createLobbyButton = button("Create a lobby")
                createLobbyButton.action {
                    serverMenu.text = "127.1.1.9"
                    controller.newLobby()
                }

                button("Join TEST").action {
                    Network.getInstance().send(GetLobbiesMessage(1))
//                    replaceWith<GameView>()
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
        controller.init(this@MainMenuView)
    }

    class JoinServerView() : View() {
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
            // Flush changes from the text fields into the model
            model.commit()

            // The edited person is contained in the model


            // A real application would persist the person here
            println("Saving ${model.hostName} / ${model.port}")
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



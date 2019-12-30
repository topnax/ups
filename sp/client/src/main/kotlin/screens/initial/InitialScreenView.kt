package screens.initial

import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import javafx.scene.control.Button
import javafx.scene.control.MenuItem
import javafx.scene.control.TextField

import tornadofx.*

class InitialScreenView : View() {


    fun setNetworkElementsEnabled(b: Boolean) {
        joinButton.disableProperty().set(!b)
    }

    lateinit var serverMenu: MenuItem
    lateinit var nameTextField: TextField
    lateinit var joinButton: Button

    val controller: InitialScreenController by inject()

    override val root = borderpane {
        top = menubar {
            serverMenu = menu("127.0.0.1") {
                item("Change server") {
                    action {
                        val modal = find<JoinServerView>()
                        modal.initialScreenController = controller
                        modal.openModal()
                    }
                }
            }
        }

        center = vbox(spacing = 10.0) {
            padding = insets(10)
            prefWidth = 10.0
            hbox(spacing = 10.0) {

                nameTextField = textfield {}

                joinButton = button("Join")
                joinButton.disableProperty().set(false)

                joinButton.action {
                    controller.onJoinButtonPressed()
                }

            }
        }
        controller.init(this@InitialScreenView)
    }

    override fun onUndock() {
        super.onUndock()
        controller.onUndock()
    }

    class JoinServerView: View() {
        var initialScreenController: InitialScreenController? = null
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

            initialScreenController?.connectTo(model.hostName.value, model.port.value.toInt())
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



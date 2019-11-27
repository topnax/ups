package views.initialscreen

import javafx.application.Platform
import javafx.scene.control.*
import mu.KotlinLogging
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*

import tornadofx.*

val logger = KotlinLogging.logger { }

class InitialScreenController : Controller(), ConnectionStatusListener {

    lateinit var initialScreen: InitialScreen

    fun init(mainMenuView: InitialScreen) {
        this.initialScreen = mainMenuView
        mainMenuView.primaryStage.setOnCloseRequest {
            alert(Alert.AlertType.INFORMATION, "Primary stage closing")
            logger.debug { "Primary stage closing" }
            Network.getInstance().send(UserLeavingMessage(), callAfterWrite = { Network.getInstance().stop() })
        }
        Network.getInstance().connectionStatusListeners.add(this)
        connectTo("localhost", 10000)
    }

    override fun onConnected() {
        Platform.runLater {
            initialScreen.setNetworkElementsEnabled(true)
            initialScreen.serverMenu.text = "Connected to ${Network.getInstance().tcpLayer?.hostname}"
        }
    }

    override fun onUnreachable() {
        Platform.runLater {
            initialScreen.setNetworkElementsEnabled(true)
            initialScreen.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} is unreachable"
        }
    }

    override fun onFailedAttempt(attempt: Int) {
        Platform.runLater {
            initialScreen.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} did not respond. Attempt $attempt"
        }
    }


    fun connectTo(hostname: String, port: Int) {
        initialScreen.setNetworkElementsEnabled(false)
        Network.getInstance().connectTo(hostname, port)
    }

    private fun validateName(): Boolean {
        if (initialScreen.nameTextField.text.trim().isNotEmpty()) {
            return true
        }
        alert(Alert.AlertType.ERROR, "Error", "Name must be of non-zero length")
        return false
    }

    fun onJoinButtonPressed() {
        if (validateName()) {
            Network.getInstance().send(UserAuthenticationMessage(initialScreen.nameTextField.text), { am: ApplicationMessage ->
                Platform.runLater {
                    if (am is UserAuthenticatedResponse) {
//                        initialScreen.replaceWith<MainMen>()
                    }
                }
            })
        }
    }
}



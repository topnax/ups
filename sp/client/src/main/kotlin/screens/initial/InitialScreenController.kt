package screens.initial

import javafx.application.Platform
import javafx.scene.control.*
import mu.KotlinLogging
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import screens.ServerRestartedEvent
import screens.ServerRestartedUnauthorizedEvent
import screens.game.GameStateRegenerationEvent
import screens.game.GameView
import screens.mainmenu.MainMenuView
import tornadofx.*

val logger = KotlinLogging.logger { }

class InitialScreenController : Controller(), ConnectionStatusListener {

    lateinit var initialScreenView: InitialScreenView

    init {
    }

    fun init(mainMenuView: InitialScreenView) {
        this.initialScreenView = mainMenuView
        mainMenuView.primaryStage.setOnCloseRequest {
            alert(Alert.AlertType.INFORMATION, "Primary stage closing")
            logger.debug { "Primary stage closing" }
            Network.getInstance().send(UserLeavingMessage(), callAfterWrite = { Network.getInstance().stop() })
        }
        Network.getInstance().connectionStatusListeners.add(this)
        connectTo("localhost", 10000)
    }

    internal fun onGameStateRegeneration(message: GameStateRegenerationResponse) {
        initialScreenView.replaceWith<GameView>()
        fire(GameStateRegenerationEvent(message))
    }

    override fun onConnected() {
        Platform.runLater {
            initialScreenView.setNetworkElementsEnabled(true)
            initialScreenView.serverMenu.disableProperty().set(false)
            initialScreenView.serverMenu.text = "Connected to ${Network.getInstance().tcpLayer?.hostname}"
        }
    }

    override fun onUnreachable() {
        Platform.runLater {
            initialScreenView.setNetworkElementsEnabled(false)
            initialScreenView.serverMenu.disableProperty().set(false)
            initialScreenView.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} is unreachable"
        }
    }

    override fun onFailedAttempt(attempt: Int) {
        Platform.runLater {
            initialScreenView.setNetworkElementsEnabled(false)
            initialScreenView.serverMenu.disableProperty().set(true)
            initialScreenView.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} did not respond. Attempt $attempt"
        }
    }

    override fun onReconnected() {

    }

    fun connectTo(hostname: String, port: Int) {
        initialScreenView.setNetworkElementsEnabled(false)
        Network.getInstance().connectTo(hostname, port)
    }

    private fun validateName(): Boolean {
        if (initialScreenView.nameTextField.text.trim().isNotEmpty()) {
            return true
        }
        alert(Alert.AlertType.ERROR, "Error", "Name must be of non-zero length")
        return false
    }

    fun onJoinButtonPressed() {
        if (validateName()) {
            Network.getInstance().send(UserAuthenticationMessage(initialScreenView.nameTextField.text), { am: ApplicationMessage ->
                Platform.runLater {
                    if (am is UserAuthenticatedResponse) {
                        Network.User = am.user
                        Network.authorized = true
                        initialScreenView.replaceWith<MainMenuView>()

                    } else if (am is GameStateRegenerationResponse) {
                        Network.User = am.user
                        Network.authorized = true
                        onGameStateRegeneration(am)
                    }
                }
            })
        }
    }

    fun onUndock() {
        Network.getInstance().removeMessageListener(::onGameStateRegeneration)
    }
}



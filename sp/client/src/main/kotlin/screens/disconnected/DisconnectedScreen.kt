package screens.disconnected

import javafx.application.Platform
import javafx.geometry.Pos
import javafx.scene.control.Alert
import mu.KotlinLogging
import screens.*
import screens.game.GameStateRegenerationEvent
import screens.game.GameView
import screens.initial.InitialScreenView
import screens.mainmenu.MainMenuView
import tornadofx.*

val logger = KotlinLogging.logger { }

class DisconnectedScreenView : View() {

    init {
        subscribe<ServerRestartedEvent> {
            Platform.runLater {
                replaceWith<MainMenuView>()
                alert(Alert.AlertType.WARNING, "Server restarted and you have been reconnected...")
            }
        }

        subscribe<MovedToLobbyScreenEvent> {
            Platform.runLater {
                replaceWith<MainMenuView>()
                alert(Alert.AlertType.WARNING, "You have been moved out of the lobby...")
            }
        }

        subscribe<NothingHappenedEvent> {
            Platform.runLater {
                replaceWith<MainMenuView>()
                alert(Alert.AlertType.INFORMATION, "Nothing happened...")
            }
        }

        subscribe<ServerRestartedUnauthorizedEvent> {
            Platform.runLater { replaceWith<InitialScreenView>() }
            alert(Alert.AlertType.WARNING, "Server restarted and user of your name is already logged in...")
        }

        subscribe<ServerUnreachableEvent> {
            Platform.runLater {
                replaceWith<InitialScreenView>()
                alert(Alert.AlertType.ERROR, "Server is unreachable...")
            }
        }

        subscribe<GameRegeneratedEvent> {
            Platform.runLater {
                replaceWith<GameView>()
                logger.warn { "Firing a stage regeneration event" }
                fire(GameStateRegenerationEvent(it.response))
            }
        }

    }

    override val root = vbox(spacing = 10.0) {
        label("Disconnected... Trying to reconnect")
        alignment = Pos.CENTER
    }


}



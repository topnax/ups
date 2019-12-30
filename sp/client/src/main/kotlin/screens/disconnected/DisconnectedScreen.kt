package screens.disconnected

import javafx.application.Platform
import javafx.geometry.Pos
import mu.KotlinLogging
import screens.GameRegeneratedEvent
import screens.ServerRestartedEvent
import screens.ServerRestartedUnauthorizedEvent
import screens.ServerUnreachableEvent
import screens.game.GameStateRegenerationEvent
import screens.game.GameView
import screens.initial.InitialScreenView
import screens.mainmenu.MainMenuView
import tornadofx.View
import tornadofx.borderpane
import tornadofx.label
import tornadofx.vbox

val logger = KotlinLogging.logger { }

class DisconnectedScreenView : View() {

    init {
        subscribe<ServerRestartedEvent> {
            Platform.runLater { replaceWith<MainMenuView>() }
        }

        subscribe<ServerRestartedUnauthorizedEvent> {
            Platform.runLater { replaceWith<InitialScreenView>() }
        }

        subscribe<ServerUnreachableEvent> {
            Platform.runLater { replaceWith<InitialScreenView>() }
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



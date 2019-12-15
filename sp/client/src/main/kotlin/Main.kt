import screens.game.GameView
import tornadofx.App
import tornadofx.launch
import screens.initialscreen.InitialScreen

class KrisKrosApp : App(InitialScreen::class)

fun main(args: Array<String>) {
    launch<KrisKrosApp>(args)
}
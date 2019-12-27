import tornadofx.App
import tornadofx.launch
import screens.initial.InitialScreen

class KrisKrosApp : App(InitialScreen::class)

fun main(args: Array<String>) {
    launch<KrisKrosApp>(args)
}
import screens.initial.InitialScreenView
import tornadofx.App
import tornadofx.launch

class KrisKrosApp : App(InitialScreenView::class)

fun main(args: Array<String>) {
    launch<KrisKrosApp>(args)
}
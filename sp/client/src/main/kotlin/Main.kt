import mu.KotlinLogging
import tornadofx.App
import tornadofx.launch
import screens.initial.InitialScreenView

class KrisKrosApp : App(InitialScreenView::class)

fun main(args: Array<String>) {
    launch<KrisKrosApp>(args)
}
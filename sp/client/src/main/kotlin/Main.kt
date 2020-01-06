import javafx.stage.Stage
import screens.initial.InitialScreenView
import tornadofx.App
import tornadofx.launch

class KrisKrosApp : App(InitialScreenView::class) {
    override fun start(stage: Stage) {
        with(stage) {
            minWidth = 575.0
            minHeight = 875.0
            super.start(this)
        }
    }
}

fun main(args: Array<String>) {
    launch<KrisKrosApp>(args)
}
import tornadofx.*

class KrisKrosApp: App(MainMenuView::class)

fun main(args: Array<String>) {
    launch<KrisKrosApp>(args)
}
package controller

import javafx.collections.ObservableList
import model.lobby.Lobby
import tornadofx.Controller
import tornadofx.observableList
import java.util.*

class MainMenuController : Controller() {

    val random = Random()

    fun newLobby() {
        lobbies.add(Lobby(1, "magor", random.nextInt(102)))
    }

    var lobbies: ObservableList<Lobby> = observableList()
}
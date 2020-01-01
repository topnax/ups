package model.lobby

import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import tornadofx.getValue
import tornadofx.setValue

/**
 * Loby ViewModel
 */
class LobbyViewModel(id: Int, owner: String, playerCount: Int) {

    val idProperty = SimpleIntegerProperty(id)
    var id by idProperty

    val ownerProperty = SimpleStringProperty(owner)
    var ownerName by ownerProperty

    val playersProperty = SimpleStringProperty("$playerCount / 4")
    var playerCount by playersProperty

}

/**
 * A class representing a lobby
 */
class Lobby(val players: List<Player>, val id: Int, val owner: Player)

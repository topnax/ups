package model.lobby

import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import tornadofx.getValue
import tornadofx.setValue

class LobbyViewModel(id: Int, owner: String, playerCount: Int) {

    val idProperty = SimpleIntegerProperty(id)
    var id by idProperty

    val ownerProperty = SimpleStringProperty(owner)
    var ownerName by ownerProperty

    val playersProperty = SimpleIntegerProperty(playerCount)
    var playerCount by playersProperty

}

class Lobby(val players: List<Player>, val id: Int, val owner: Player)
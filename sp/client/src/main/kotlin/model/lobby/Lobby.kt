package model.lobby

import javafx.beans.property.SimpleIntegerProperty
import javafx.beans.property.SimpleStringProperty
import javafx.collections.ObservableList
import tornadofx.ViewModel
import tornadofx.getValue
import tornadofx.observableList
import tornadofx.setValue

class Lobby(id: Int, owner: String, players: Int) {

    val idProperty = SimpleIntegerProperty(id)
    var id by idProperty

    val ownerProperty = SimpleStringProperty(owner)
    var owner by ownerProperty

    val playersProperty = SimpleIntegerProperty(players)
    var players by playersProperty

}
package controller

import javafx.collections.ObservableList
import model.lobby.Lobby
import MainMenuView
import javafx.application.Platform
import networking.*
import networking.receiver.Message
import networking.reader.MessageReader
import networking.receiver.SimpleMessageReceiver
import tornadofx.Controller
import tornadofx.observableList
import java.util.*

class MainMenuController : Controller(), MessageReader, ConnectionStatusListener {

    val random = Random()

    lateinit var mainMenuView: MainMenuView

    lateinit var tcp: TCPLayer

    var lobbies: ObservableList<Lobby> = observableList()

    override fun onConnected() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
            mainMenuView.serverMenu.text = "Connected to ${tcp.hostname}"
        }
    }

    override fun onUnreachable() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
            mainMenuView.serverMenu.text = "${tcp.hostname} is unreachable"
        }
    }

    override fun onFailedAttempt(attempt: Int) {
        Platform.runLater {
            mainMenuView.serverMenu.text = "${tcp.hostname} did not respond. Attempt $attempt"
        }
    }

    @Synchronized
    override fun read(message: Message) {
        mainMenuView.serverMenu.text = message.content
    }

    fun newLobby() {
        lobbies.add(Lobby(1, "magor", random.nextInt(102)))
    }

    fun init(mainMenuView: MainMenuView) {
        this.mainMenuView = mainMenuView
        mainMenuView.setNetworkElementsEnabled(false)
        this.tcp = TCPLayer(
                connectionStatusListener = this,
                messageReceiver = SimpleMessageReceiver(object : MessageReader {
                    override fun read(message: Message) {

                    }
                })
        )
        tcp.start()
    }

}

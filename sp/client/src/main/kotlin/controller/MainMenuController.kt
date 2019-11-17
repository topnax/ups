package controller

import javafx.collections.ObservableList
import model.lobby.Lobby
import MainMenuView
import javafx.application.Platform
import networking.*
import networking.messages.ApplicationMessage
import networking.messages.GetLobbiesMessage
import networking.messages.LobbiesListMessage
import networking.receiver.Message
import networking.reader.MessageReader
import networking.receiver.SimpleMessageReceiver
import tornadofx.Controller
import tornadofx.observableList
import java.util.*

class MainMenuController : Controller(), MessageReader, ConnectionStatusListener {

    val random = Random()

    lateinit var mainMenuView: MainMenuView

    var lobbies: ObservableList<Lobby> = observableList()

    init {
        Network.getInstance().addMessageListener(LobbiesListMessage::class.java, ::onLobbiesListRetrieved)
    }

    fun onLobbiesListRetrieved(message: ApplicationMessage) {
        if (message is LobbiesListMessage) {
            message.lobbies.forEach{
                println("lobby ${it.id}")
                println(it.owner)
                println(it.id)
                println(it.players)
            }
        }
    }

    override fun onConnected() {
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
            mainMenuView.serverMenu.text = "Connected to ${Network.getInstance().tcpLayer?.hostname}"
            Network.getInstance().send(GetLobbiesMessage(1))
        }
    }

    override fun onUnreachable() {
        println("onunreachable")
        Platform.runLater {
            mainMenuView.setNetworkElementsEnabled(true)
            mainMenuView.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} is unreachable"
        }
    }

    override fun onFailedAttempt(attempt: Int) {
        println("onfailed" + attempt)
        Platform.runLater {
            mainMenuView.serverMenu.text = "${Network.getInstance().tcpLayer?.hostname} did not respond. Attempt $attempt"
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
        Network.getInstance().connectionStatusListeners.add(this)
        connectTo("localhost", 10000)
    }

    fun connectTo(hostname: String, port: Int) {
        mainMenuView.setNetworkElementsEnabled(false)
        Network.getInstance().connectTo(hostname, port)
    }

}

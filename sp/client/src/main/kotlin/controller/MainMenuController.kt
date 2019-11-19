package controller

import MainMenuView
import javafx.application.Platform
import javafx.collections.ObservableList
import model.lobby.Lobby
import networking.ConnectionStatusListener
import networking.Network
import networking.messages.*
import networking.reader.MessageReader
import networking.receiver.Message
import tornadofx.Controller
import tornadofx.observableList
import views.LobbyView
import java.util.*

class MainMenuController : Controller(), MessageReader, ConnectionStatusListener {

    val random = Random()

    lateinit var mainMenuView: MainMenuView

    var lobbies: ObservableList<Lobby> = observableList()

    init {
        Network.getInstance().addMessageListener(GetLobbiesMessage::class.java, ::onLobbiesListRetrieved)
    }

    fun onLobbiesListRetrieved(message: ApplicationMessage) {
        if (message is GetLobbiesMessage) {
            message.lobbies.forEach {
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
//            Network.getInstance().send(GetLobb(1))
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

    fun newLobby(name: String) {
//        lobbies.add(Lobby(1, "magor", random.nextInt(102)))

        Network.getInstance().send(CreateLobbyMessage(name), { am: ApplicationMessage ->
            run {
                when (am) {
                    is SuccessResponseMessage -> {
                        println("created yes ${am.content}")
                        Platform.runLater {
                            mainMenuView.replaceWith(find<LobbyView>(mapOf(LobbyView::playerName to mainMenuView.nameTextField.text)))
                        }

                    }
                    is ErrorResponseMessage -> println("fail ${am.content}")
                    else -> {

                    }
                }
            }
        })
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

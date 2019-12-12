package screens.game

import model.game.Desk
import model.game.Letter
import model.game.Tile
import model.lobby.Player
import mu.KotlinLogging
import tornadofx.Controller
import tornadofx.EventBus
import tornadofx.FXEvent
import java.util.*
import kotlin.concurrent.schedule

private val logger = KotlinLogging.logger { }

class GameScreenController : Controller() {

    var selectedTile: Tile? = null

    val desk = Desk(Desk.getTilesFromTileTypes())

    var letters: MutableList<Letter> = mutableListOf()

    val placedLetters: MutableList<Letter> = mutableListOf()

    lateinit var gameView: GameView

    var players: List<Player> = emptyList()

    var activePlayerID = -1

    fun init(gameView: GameView) {
        this.gameView = gameView
        Timer().schedule(5000){
            fire(
                    NewLetterSackEvent(
                            listOf(
                                    Letter("CH", 2),
                                    Letter("Q", 2),
                                    Letter("V", 2),
                                    Letter("B", 2),
                                    Letter("S", 2)
                            )
                    )
            )
            players = listOf(
                    Player("Standa", 1, false),
                    Player("Pavel", 2, false),
                    Player("Lobotom", 3, false)
            )
            activePlayerID = 2
            fire(PlayerStateChangedEvent())
        }

    }

    init {
        subscribe<TileSelectedEvent> {
            logger.debug { "Tile $it.tile selected" }
            selectedTile = it.tile
        }

        subscribe<LetterPlacedEvent> {
            logger.debug { "Letter of value ${it.letter} has been pressed GO!" }
            val tile = desk.tiles[selectedTile!!.row][selectedTile!!.column]
            tile.letter = it.letter
            tile.letter?.run {
                letters.remove(it.letter)
                placedLetters.add(it.letter)
                tile.selected = false
                fire(DeskChange(tile))
            }
        }

        subscribe<NewLetterSackEvent> {
            letters = it.letters.toMutableList()
        }

        subscribe<TileWithLetterClicked> {
            if (placedLetters.contains(it.tile.letter)) {
                placedLetters.remove(it.tile.letter)
                letters.add(it.tile.letter!!)
                desk.tiles[it.tile.row][it.tile.column].letter = null
                fire(DeskChange(it.tile))
                fire(NewLetterSackEvent(letters))
            }
        }
    }
}

class NewLetterSackEvent(val letters: List<Letter>) : FXEvent(EventBus.RunOn.BackgroundThread)

class DeskChange(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class TileWithLetterClicked(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class LetterPlacedEvent(val letter: Letter) : FXEvent(EventBus.RunOn.BackgroundThread)

class TileSelectedEvent(val tile: Tile) : FXEvent(EventBus.RunOn.BackgroundThread)

class PlayerStateChangedEvent() : FXEvent(EventBus.RunOn.BackgroundThread)


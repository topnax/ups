package screens.game

import javafx.application.Platform
import javafx.geometry.Insets
import javafx.geometry.Pos
import javafx.scene.control.Button
import javafx.scene.layout.GridPane
import javafx.scene.paint.Color
import javafx.scene.text.FontWeight
import model.game.Letter
import model.game.Tile
import model.game.TileType
import model.lobby.Player
import mu.KotlinLogging
import tornadofx.*

private val logger = KotlinLogging.logger { }

class GameView : View() {

    val controller: GameScreenController by inject()

    val tileViews = arrayOfNulls<Array<TileView?>>(15)

    init {
        controller.init(this)
    }

    override val root = vbox {
        println(controller)
        gridpane {

            style {
                backgroundColor += Color.PINK
            }
            useMaxWidth = true
            var pindex = 0
            controller.desk.tiles.forEachIndexed { index, tile ->
                tileViews[index] = arrayOfNulls<TileView>(15)
                tile.forEachIndexed { i, tile ->

                    if (pindex % 3 == 0) {
                        tile.letter = Letter("s", 30)
                    }

                    if (pindex == 21) {
                        tile.letter = Letter("ch", 30)
                    }

                    pindex++

                    val tv = TileView(tile)
                    this@gridpane.add(tv)
                    tileViews[index]!![i] = tv
                }
            }
            gridLinesVisibleProperty().set(true)
            subscribe<DeskChange> {
                refreshTile(it.tile)
            }

            subscribe<TileSelectedEvent> {
                logger.debug { "Tile event" }
                controller.selectedTile?.let {
                    it.selected = false
                    refreshTile(it)
                }
                if (it.tile.letter == null) {
                    it.tile.selected = true
                    refreshTile(it.tile)
                }
            }

        }
        hbox(spacing = 10) {
            subscribe<PlayerStateChangedEvent> {
                logger.debug { "player state changed event" }
                Platform.runLater {
                    clear()
                    controller.players.forEach { label(it.name) { if (it.id == controller.activePlayerID) style { fontWeight = FontWeight.EXTRA_BOLD } } }
                }
            }
        }
        tilepane {
            hgap = 10.0
            vgap = 10.0
            prefColumns = 4

            subscribe<NewLetterSackEvent> { event ->
                Platform.runLater {
                    clear()
                    event.letters.forEach { add(LetterView(it)) }
                }
            }
        }


    }

    fun GridPane.refreshTile(tile: Tile) {
        val tv = TileView(tile)
        Platform.runLater {
            tileViews[tile.column]!![tile.row]?.removeFromParent()
            this.add(tv)
        }
        tileViews[tile.column]!![tile.row] = tv
    }
}

class LetterView(val letter: Letter) : View() {

    override val root = hbox() {
        padding = Insets(5.0)
        alignment = Pos.BOTTOM_CENTER
        label(letter.value) {
            textFill = Color.WHITE
            style {
                fontSize = 16.px
            }
        }
        label(letter.points.toString()) {
            textFill = Color.WHITE
            style {
                fontSize = 12.px
            }
        }
        style {
            backgroundColor += Color.DARKGREEN
        }
        setOnMouseClicked {
            removeFromParent()
            fire(LetterPlacedEvent(letter))
            logger.debug { "Letter ${letter.value} clicked" }
        }
    }
}

class TileView(val tile: Tile) : View() {
    lateinit var btn: Button

    override val root = tile.letter?.let {
        hbox {
            padding = Insets(6.0)
            alignment = Pos.CENTER
            gridpaneConstraints {
                columnRowIndex(tile.column, tile.row)
            }
            btn = button(it.value.toUpperCase()) {
                action {
                    fire(TileWithLetterClicked(tile))
                }
                textFill = Color.WHITE
                style {
                    backgroundColor += Color.GREEN

                    borderRadius += box(6.px)
                }

            }
            style {
                backgroundColor += if (!tile.selected) tile.type.getTileColor() else Color.PINK
            }
            btn
        }

    } ?: run {
        button {
            gridpaneConstraints {
                fillHeightWidth = true
                useMaxSize = true
                columnRowIndex(tile.column, tile.row)
            }
            action {
                logger.info { "Pressed at ${tile.row}#${tile.column}" }
                fire(TileSelectedEvent(tile))
            }
            style {
                backgroundColor += if (!tile.selected) tile.type.getTileColor() else Color.PINK
            }
        }
    }
}

fun TileType.getTileColor(): Color {
    return when (this) {
        TileType.BASIC -> Color.LIGHTYELLOW
        TileType.MULTIPLY_LETTER_2 -> Color.LIGHTSALMON
        TileType.MULTIPLY_LETTER_3 -> Color.DARKRED
        TileType.MULTIPLY_WORD_2 -> Color.LIGHTGREEN
        TileType.MULTIPLY_WORD_3 -> Color.DARKGREEN
    }
}
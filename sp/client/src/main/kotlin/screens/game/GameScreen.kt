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
import mu.KotlinLogging
import networking.Network
import tornadofx.*

private val logger = KotlinLogging.logger { }

class GameView : View() {

    val controller: GameScreenController by inject()

    val tileViews = arrayOfNulls<Array<TileView?>>(15)

    lateinit var finishButton: Button

    lateinit var declineWordsButton: Button

    lateinit var acceptWordsButton: Button

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
            controller.desk.tiles.forEachIndexed { index, row ->
                tileViews[index] = arrayOfNulls<TileView>(15)
                row.forEachIndexed { i, tile ->
                    val tv = TileView(tile)
                    this@gridpane.add(tv)
                    tileViews[index]!![i] = tv
                }
            }
            gridLinesVisibleProperty().set(true)
            subscribe<DeskChange> {
                Platform.runLater {
                    refreshTile(it.tile)
                }
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
                logger.debug { "Player state changed event" }
                Platform.runLater {
                    clear()
                    controller.players.forEach {
                        label(it.name) {
                            if (it.id == controller.activePlayerID) style { fontWeight = FontWeight.EXTRA_BOLD }
                        }
                        if (controller.playerIdsWhoAcceptedWords.contains(it.id)) label("Accepted words!!!")
                    }
                    label(if (controller.activePlayerID == Network.User.id) "Jste na tahu" else "Nejste na tahu")

                    logger.debug {"controller.activePlayerID != Network.User.id: ${controller.activePlayerID != Network.User.id}"}
                    logger.debug {"controller.roundFinished: ${controller.roundFinished}"}
                    logger.debug {"!controller.wordsAccepted: ${!controller.wordsAccepted}"}

                    acceptWordsButton.visibleProperty().set(controller.activePlayerID != Network.User.id && controller.roundFinished && !controller.wordsAccepted)
                    declineWordsButton.visibleProperty().set(controller.activePlayerID != Network.User.id && controller.roundFinished && !controller.wordsAccepted)
                    finishButton.visibleProperty().set(!controller.roundFinished && controller.activePlayerID == Network.User.id)

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

        finishButton = button("Finish round") {
            subscribe<GameStartedEvent> {
                visibleProperty().set(it.message.activePlayerId == Network.User.id)
                action {
                    controller.onFinishRoundButtonClicked()
                }
            }
        }

        hbox (spacing=10) {
            acceptWordsButton = button("Accept words") {
                visibleProperty().set(false)
                action {
                    controller.onAcceptWordsButtonClicked()
                }
                subscribe<RoundFinishedEvent> {
                    logger.debug {"Round finished event: controller.activePlayerID != Network.User.id ${controller.activePlayerID != Network.User.id}"}
                    visibleProperty().set(controller.activePlayerID != Network.User.id)
                }
            }

            declineWordsButton = button("Decline words") {
                visibleProperty().set(false)
                action {
                    controller.onDeclineWordsClicked()
                }
                subscribe<RoundFinishedEvent> {
                    visibleProperty().set(controller.activePlayerID != Network.User.id)
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

    override fun onDock() {
        super.onDock()
        controller.onDock()
    }

    override fun onUndock() {
        super.onUndock()
        controller.onUndock()
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
            fire(LetterPlacedEvent(letter, this@LetterView))
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
                    backgroundColor += if (tile.highlighted) Color.ORANGE else Color.GREEN

                    borderRadius += box(6.px)
                }

            }
            style {
                backgroundColor += if (!tile.selected) tile.typeEnum.getTileColor() else Color.PINK
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
                backgroundColor += if (!tile.selected) tile.typeEnum.getTileColor() else Color.PINK
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
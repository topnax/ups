package screens.game

import javafx.application.Platform
import javafx.geometry.Insets
import javafx.geometry.Pos
import javafx.scene.control.Button
import javafx.scene.layout.GridPane
import javafx.scene.layout.Priority
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

    val affectedTileCoordinates = mutableSetOf<Pair<Int, Int>>()

    init {
        controller.init(this)
    }

    override val root = vbox(spacing = 10) {
        padding = Insets(10.0)
        label(Network.User.name)
        gridpane {
            useMaxWidth = true
            hgrow = Priority.ALWAYS
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
                    gridLinesVisibleProperty().set(true)
                }
            }
        }
        vbox(spacing = 10) {
            subscribe<PlayerStateChangedEvent> {
                logger.debug { "Player state changed event" }
                Platform.runLater {
                    clear()
                    controller.players.forEach {
                        label("${it.name} (${controller.playerPointsMap[it.id]})" + if (controller.activePlayerID == it.id) " <${controller.currentRoundPlayerPoints}>" else "") {
                            if (it.id == controller.activePlayerID) style { fontWeight = FontWeight.EXTRA_BOLD }
                        }
                        if (it.disconnected) label("disconnected") { style { fontWeight = FontWeight.EXTRA_BOLD } }
                        if (controller.playerIdsWhoAcceptedWords.contains(it.id)) label("Accepted words!!!")
                    }
                    label(if (controller.activePlayerID == Network.User.id) "It's your turn" else "It is not your turn") {
                        style {
                            fontWeight = FontWeight.EXTRA_BOLD
                        }
                    }

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
            visibleProperty().set(controller.activePlayerID == Network.User.id)
            action {
                controller.onFinishRoundButtonClicked()
            }
        }

        hbox(spacing = 10) {
            acceptWordsButton = button("Accept words") {
                visibleProperty().set(false)
                action {
                    controller.onAcceptWordsButtonClicked()
                }
                subscribe<RoundFinishedEvent> {
                    logger.debug { "Round finished event: controller.activePlayerID != Network.User.id ${controller.activePlayerID != Network.User.id}" }
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
        logger.info { "Refreshing at c${tile.column}#r${tile.row} " }
        affectedTileCoordinates.add(Pair(tile.column, tile.row))
    }

    override fun onDock() {
        super.onDock()
        resetDesk()
        controller.onDock()
    }

    override fun onUndock() {
        super.onUndock()
        controller.onUndock()
    }

    private fun resetDesk() {
        logger.info { "About to reset ${affectedTileCoordinates.size} tiles" }
        affectedTileCoordinates.forEach {
            logger.info { "Reseting at c${it.first}#r${it.second}" }
            controller.desk.tiles[it.second][it.first].letter = null
            controller.desk.tiles[it.second][it.first].set = false
            fire(DeskChange(controller.desk.tiles[it.second][it.first]))
        }
        affectedTileCoordinates.clear()
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
            usePrefSize = true
            prefHeight(60.0)
            prefWidth(60.0)
            minWidth(100.0)
            minHeight(150.0)
            padding = Insets(5.0)

            alignment = Pos.CENTER
            gridpaneConstraints {
                columnRowIndex(tile.column, tile.row)
            }

            hbox {
                padding = Insets(8.0)
                alignment = Pos.BOTTOM_RIGHT
                label(it.value.toUpperCase()) {
                    textFill = Color.WHITE
//                    setOnMouseClicked {
//                        fire(TileWithLetterClicked(tile))
//                    }
                }
                label(it.points.toString()) {
                    textFill = Color.WHITE
                    style {
                        fontSize = 10.px
                    }
//                    setOnMouseClicked {
//                        fire(TileWithLetterClicked(tile))
//                    }
                }
                style {
                    backgroundColor += if (tile.highlighted) Color.ORANGE else Color.GREEN
                    borderRadius += box(6.px)
                }

                setOnMouseClicked {
                    fire(TileWithLetterClicked(tile))
                }
            }

//            btn = button(it.value.toUpperCase()) {
//                action {
//                    fire(TileWithLetterClicked(tile))
//                }
//                textFill = Color.WHITE
//                style {
//                    backgroundColor += if (tile.highlighted) Color.ORANGE else Color.GREEN
//
//                    borderRadius += box(6.px)
//                }
//
//            }
            style {
                backgroundColor += if (!tile.selected) tile.typeEnum.getTileColor() else Color.PINK
            }
//            btn
        }
    } ?: run {
        hbox {
            usePrefSize = true
            prefHeight(60.0)
            prefWidth(60.0)
            minWidth(100.0)
            minHeight(150.0)
            padding = Insets(5.0)

            alignment = Pos.CENTER
            gridpaneConstraints {
                columnRowIndex(tile.column, tile.row)
            }

            hbox {
                padding = Insets(8.0)
                alignment = Pos.BOTTOM_RIGHT
                label(" ") {
                    textFill = Color.WHITE
                }
                label(" ") {
                    textFill = Color.WHITE
                    style {
                        fontSize = 10.px
                    }
                }


                setOnMouseClicked {
                    logger.info { "Pressed at ${tile.row}#${tile.column}" }
                    fire(TileSelectedEvent(tile))
                }
            }

//            btn = button(it.value.toUpperCase()) {
//                action {
//                    fire(TileWithLetterClicked(tile))
//                }
//                textFill = Color.WHITE
//                style {
//                    backgroundColor += if (tile.highlighted) Color.ORANGE else Color.GREEN
//
//                    borderRadius += box(6.px)
//                }
//
//            }
            style {
                backgroundColor += if (!tile.selected) tile.typeEnum.getTileColor() else Color.PINK
            }
//            btn
        }
//        hbox {
//            usePrefSize = true
//            prefHeight(100.0)
//            prefWidth(100.0)
//            prefHeight(150.0)
//            prefWidth(150.0)
//            minWidth(100.0)
//            minHeight(150.0)
//            gridpaneColumnConstraints {
//                prefHeight(150.0)
//                prefWidth(150.0)
//                prefHeight(200.0)
//                prefWidth(200.0)
//                minWidth(100.0)
//                minHeight(150.0)
//            }
//            gridpaneConstraints {
//                fillHeightWidth = true
//                useMaxSize = true
//                columnRowIndex(tile.column, tile.row)
//            }
//            setOnMouseClicked {
//                logger.info { "Pressed at ${tile.row}#${tile.column}" }
//                fire(TileSelectedEvent(tile))
//            }
//            style {
//                backgroundColor += if (!tile.selected) tile.typeEnum.getTileColor() else Color.PINK
//            }
//            label("C") {
//                textFill = Color.WHITE
//                style {
//                    fontSize = 10.px
//                }
//            }
//        }
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
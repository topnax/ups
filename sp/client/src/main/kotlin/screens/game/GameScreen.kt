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
import screens.UserAuthenticatedEvent
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
        prefWidth = 800.0
        prefHeight = 600.0
        padding = Insets(10.0)
        hbox(spacing = 10.0) {
            val nameLabel = label(Network.User.name)

            subscribe<UserAuthenticatedEvent> {
                nameLabel.text = it.name
            }

            button("Leave the game") {
                action {
                    controller.leaveGame()
                }
            }

        }
        gridpane {
            gridLinesVisibleProperty().set(true)
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
                        hbox {
                            label("${it.name} (${controller.playerPointsMap[it.id]})" + if (controller.activePlayerID == it.id) " <${controller.currentRoundPlayerPoints}>" else "") {
                                if (it.id == controller.activePlayerID) style { fontWeight = FontWeight.EXTRA_BOLD }
                            }
                            if (it.disconnected) label("disconnected") { style { fontWeight = FontWeight.EXTRA_BOLD } }
                            if (controller.playerIdsWhoAcceptedWords.contains(it.id)) label("Accepted words")
                        }
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
            subscribe<PlayerStateChangedEvent> {
                visibleProperty().set(controller.activePlayerID == Network.User.id)
            }
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

    private fun GridPane.refreshTile(tile: Tile) {
        val tv = TileView(tile)
        Platform.runLater {
            tileViews[tile.column]!![tile.row]?.removeFromParent()
            this.add(tv)
        }
        tileViews[tile.column]!![tile.row] = tv
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

    override val root = tile.letter?.let {
        // letter layout
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
                }
                label(it.points.toString()) {
                    textFill = Color.WHITE
                    style {
                        fontSize = 10.px
                    }
                }
                style {
                    backgroundColor += if (tile.highlighted) Color.ORANGE else Color.GREEN
                    borderRadius += box(6.px)
                }

                setOnMouseClicked {
                    fire(TileWithLetterClicked(tile))
                }
            }

            style {
                backgroundColor += if (!tile.selected) tile.typeEnum.getTileColor() else Color.PINK
            }
        }
    } ?: run {
        // empty tile layout
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
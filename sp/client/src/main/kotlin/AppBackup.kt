import javafx.beans.property.ObjectProperty
import javafx.geometry.Insets
import javafx.scene.Parent
import javafx.scene.control.Button
import javafx.scene.paint.Color
import tornadofx.*
import java.util.*

class GameView : View() {
    val buttons = mutableListOf<Button>()
    var index = 0
    override val root = vbox {
        val random = Random()
        button("Back").action {
            replaceWith<MainMenuView>()
        }
        gridpane {
            id = "gpane"
            useMaxWidth = true
            for (c in 1..(15 * 15)) {
                val tv = Tile(index % 15, index / 15, TileType.BASIC)
                tv.letter = Letter("S", 10)
                this.add(TileView(tv))
                index++
            }
            style {
                backgroundColor += Color.LIGHTYELLOW
            }
            gridLinesVisibleProperty().set(true)
        }
    }

    init {
        root.lookup("#gpane")
    }


}

class TileView(tile: Tile) : View() {

    override val root = tile.letter?.let {
        button(it.value.toUpperCase()) {
            action {
                println("${it.value} has been pressed")
            }
            gridpaneConstraints {
                columnRowIndex(tile.column, tile.row)
                margin = Insets(3.0)
            }
            textFill = Color.WHITE
            style {
                backgroundColor += Color.GREEN

                borderRadius += box(6.px)
            }
            useMaxHeight = true
            useMaxWidth = true
        }
    } ?: run {
        pane()
    }
//        }
//    override val root = button(if (tile.letter != null.) {tile..toString().toUpperCase()) {
//            action {
//                println("$c has been pressed")
//                buttons[random.nextInt(15)].text = "hello idiot"
//            }
//            gridpaneConstraints {
//                columnRowIndex(index % 15, index / 15)
//                margin = Insets(3.0)
//            }
//            textFill = Color.WHITE
//            style {
//                backgroundColor += Color.GREEN
//
//                borderRadius += box(6.px)
//            }
//            useMaxHeight = true
//            useMaxWidth = true
//        }
}


class Tile(val row: Int, val column: Int, val type: TileType) {
    var letter: Letter? = null

}

enum class TileType {
    BASIC,
    MULTIPLY_WORD_2,
    MULTIPLY_WORD_3,
    MULTIPLY_LETTER_2,
    MULTIPLY_LETTER_3
}

fun Tile.getPoints(letter: Letter): Int =
        when (this.type) {
            TileType.MULTIPLY_LETTER_2 -> letter.points * 2
            TileType.MULTIPLY_LETTER_3 -> letter.points * 3
            else -> letter.points
        }


data class Letter(val value: String, val points: Int)

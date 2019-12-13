package model.game


class Tile(val row: Int, val column: Int, type: Int, var selected: Boolean = false, var highlighted: Boolean = false) {
    val type = if (type < TileType.values().size) TileType.values()[type] else TileType.BASIC
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
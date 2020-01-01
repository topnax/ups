package model.game

import com.beust.klaxon.Json

// a class representing a KrisKros tile
class Tile(val row: Int, val column: Int, val type: Int, var selected: Boolean = false, var highlighted: Boolean = false, var letter: Letter? = null, var set: Boolean = false) {
    @Json(ignored = true)
    val typeEnum = if (type < TileType.values().size) TileType.values()[type] else TileType.BASIC
}

enum class TileType {
    BASIC,
    MULTIPLY_WORD_2,
    MULTIPLY_WORD_3,
    MULTIPLY_LETTER_2,
    MULTIPLY_LETTER_3
}

fun Tile.getPoints(letter: Letter): Int =
        when (this.typeEnum) {
            TileType.MULTIPLY_LETTER_2 -> letter.points * 2
            TileType.MULTIPLY_LETTER_3 -> letter.points * 3
            else -> letter.points
        }
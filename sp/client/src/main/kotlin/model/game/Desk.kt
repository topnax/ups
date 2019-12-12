package model.game

import kotlin.collections.mapIndexed

class Desk(val tiles: Array<Array<Tile>>) {

    companion object {
        fun getInitialTileTypes(): Array<Array<TileType>> {
            return arrayOf(
                    arrayOf(TileType.MULTIPLY_WORD_3, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_3),
                    arrayOf(TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC),
                    arrayOf(TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC),
                    arrayOf(TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2),
                    arrayOf(TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC),
                    arrayOf(TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC),
                    arrayOf(TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC),

                    arrayOf(TileType.MULTIPLY_WORD_3, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_3),

                    arrayOf(TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC),
                    arrayOf(TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC),
                    arrayOf(TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.BASIC),
                    arrayOf(TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2),
                    arrayOf(TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC),
                    arrayOf(TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_3, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC),
                    arrayOf(TileType.MULTIPLY_WORD_3, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_2, TileType.BASIC, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_LETTER_2, TileType.BASIC, TileType.BASIC, TileType.MULTIPLY_WORD_3)
            )
        }

        fun getTilesFromTileTypes(): Array<Array<Tile>> {

            return getInitialTileTypes().mapIndexed { index, arrayOfTileTypes -> arrayOfTileTypes.mapIndexed { ind, type -> Tile(index, ind, type = type.ordinal) }.toTypedArray() }.toTypedArray()
//            return arr
        }
    }
}
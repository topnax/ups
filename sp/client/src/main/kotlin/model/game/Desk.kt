package model.game

import kotlin.collections.mapIndexed

/**
 * A class representing a desk, containing it's tiles
 */
class Desk(val tiles: Array<Array<Tile>>) {
    companion object {

        /**
         * Returns a 2D array of tile types, representing the tiles on a Kris Kros desk
         */
        private fun getInitialTileTypes(): Array<Array<TileType>> {
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

        /**
         * Creates 2D array of tiles based on the implicit KrisKros tiles
         */
        fun getTilesFromTileTypes(): Array<Array<Tile>> {
            return getInitialTileTypes().mapIndexed { index, arrayOfTileTypes -> arrayOfTileTypes.mapIndexed { ind, type -> Tile(index, ind, type = type.ordinal) }.toTypedArray() }.toTypedArray()
        }
    }
}
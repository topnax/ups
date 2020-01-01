package model.lobby

/**
 * A class representing a player
 */
data class Player(val name: String, val id: Int, val ready: Boolean, var disconnected: Boolean = false)
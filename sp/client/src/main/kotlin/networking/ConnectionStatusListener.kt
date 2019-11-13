package networking

interface ConnectionStatusListener {
    fun onConnected()
    fun onUnreachable()
    fun onFailedAttempt(attempt: Int)
}
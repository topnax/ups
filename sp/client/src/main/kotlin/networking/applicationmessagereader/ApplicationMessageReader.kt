package networking.applicationmessagereader

import networking.messages.ApplicationMessage

interface ApplicationMessageReader {
    fun read(message: ApplicationMessage, mid: Int)
}
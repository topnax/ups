import com.beust.klaxon.KlaxonException
import networking.applicationmessagereader.ApplicationMessageReader
import networking.messages.ApplicationMessage
import networking.messages.JoinLobbyMessage
import networking.reader.SimpleMessageReader
import networking.receiver.Message
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.fail

class SimpleMessageReaderTest {

    @Test
    internal fun validMessageReadTest() {
        val messages = mutableListOf<ApplicationMessage>()
        val reader = SimpleMessageReader(object : ApplicationMessageReader {
            override fun read(message: ApplicationMessage, mid: Int) {
                messages.add(message)
            }
        })

        reader.read(Message(15, 1, """
           {
                "player_name": "Topnax",
                "lobby_id": 10
           }
        """, 10))

        assertEquals(1, messages.size)
        assertEquals(JoinLobbyMessage::class.java, messages[0].javaClass)
        assertEquals(10, (messages[0] as JoinLobbyMessage).lobbyId)
        assertEquals(10, (messages[0] as JoinLobbyMessage).lobbyId)
    }

    @Test
    internal fun invalidMessageReadTest() {
        val messages = mutableListOf<ApplicationMessage>()
        val reader = SimpleMessageReader(object : ApplicationMessageReader {
            override fun read(message: ApplicationMessage, mid: Int) {
                messages.add(message)
            }
        })

        reader.read(Message(15, 1, """
           fdsfdsfds f5sd4 f54ds
        """, 10))

        assertEquals(0, messages.size)
    }
}
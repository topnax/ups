import networking.receiver.Message
import networking.reader.MessageReader
import networking.receiver.SimpleMessageReceiver
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test

class SimpleMessageReceiverTest {

    private val received = mutableListOf<Message>()
    private lateinit var receiver: SimpleMessageReceiver

    @BeforeEach
    internal fun setUp() {
        received.clear()
        receiver = SimpleMessageReceiver(object : MessageReader {
            override fun read(message: Message) {
                received.add(message)
            }
        })
    }

    @Test
    internal fun receiveSimpleValidMessage() {
        val message = "$10#1#123456789X"
        receiver.receive(message.toByteArray(), message.length)
        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
    }

    @Test
    internal fun receiveSimpleValidSplitMessage() {
        receiver.receive("$10#1#123".toByteArray(), 8)
        receiver.receive("456789X".toByteArray(), 7)

        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
    }

    @Test
    internal fun receiveSimpleValidMessageWithPreviousIncorrect() {
        receiver.receive("678IDS".toByteArray(), 8)
        receiver.receive("$434á438#".toByteArray(), 7)

        receiver.receive("$10#1#123".toByteArray(), 8)
        receiver.receive("456789X".toByteArray(), 7)

        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
    }

    @Test
    internal fun receiveMultipleValidMessages() {
        receiver.receive("$10#1#123".toByteArray(), 8)
        receiver.receive("456789X".toByteArray(), 7)

        receiver.receive("$17#2#123".toByteArray(), 8)
        receiver.receive("456789X".toByteArray(), 7)
        receiver.receive("456789X".toByteArray(), 7)

        assertEquals(2, received.size)

        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)

        assertEquals(17, received[1].length)
        assertEquals("123456789X456789X", received[1].content)
        assertEquals(2, received[1].type)
    }

    @Test
    internal fun receiveSingleInvalidMessage() {
        receiver.receive("$19#1#123".toByteArray(), 8)
        receiver.receive("456789X".toByteArray(), 7)

        assertEquals(0, received.size)
    }

    @Test
    internal fun receiveMultipleInvalidMessage() {
        receiver.receive("1546f#1#123".toByteArray(), 8)
        receiver.receive("1546SDs\$dsdsd23".toByteArray(), 54)
        receiver.receive("1546f#1#12$3dsd".toByteArray(), 132)
        receiver.receive("1546f#d1465sad 46s5f1#123".toByteArray(), 12)

        assertEquals(0, received.size)
    }

}
import networking.reader.MessageReader
import networking.receiver.Message
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
    internal fun simpleReceiveTest() {
        val message = "$10#1#11#123456789X$9#1#11#123456789"
        assertEquals("$10#1#11#123456789X", String(receiver.receive(message.toByteArray(), message.toByteArray().size)[0]))
        assertEquals("$9#1#11#123456789", String(receiver.receive(message.toByteArray(), message.toByteArray().size)[1]))
    }

    @Test
    internal fun simpleReceiveTest2() {
        val message = "$10#1#11#123456789X$19$9#1#11#123456789"
        val received = receiver.receive(message.toByteArray(), message.toByteArray().size)
        assertEquals("$10#1#11#123456789X", String(received[0]))
        assertEquals("$19", String(received[1]))
        assertEquals("$9#1#11#123456789", String(received[2]))
    }

    @Test
    internal fun receiveSimpleValidMessage() {
        val message = "$10#1#11#123456789X"
        receiver.receive(message.toByteArray(), message.length)
        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)
    }

    @Test
    internal fun receiveSimpleValidSplitMessage() {
        receiver.receive("$10#1#11#123".toByteArray(), "$10#1#11#123".toByteArray().size)
        receiver.receive("456789X".toByteArray(), "456789X".toByteArray().size)

        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)
    }

    @Test
    internal fun receiveSimpleMessageSplitAtAccent() {
        val content = "alzáček"
        val contentByteArray = content.toByteArray()
        val header = "$${contentByteArray.size}#1#22#"

        // receive header
        receiver.receive(header.toByteArray(), header.toByteArray().size)
        // receive "alz'"
        receiver.receive(contentByteArray.copyOfRange(0, 4), contentByteArray.copyOfRange(0, 4).size)
        // receive "aček"
        receiver.receive(contentByteArray.copyOfRange(4, contentByteArray.size), contentByteArray.copyOfRange(4, contentByteArray.size).size)

        assertEquals(1, received.size)
        assertEquals("alzáček", received[0].content)
        assertEquals("alzáček".toByteArray().size, received[0].length)
    }

    @Test
    internal fun receiveSimpleValidMessageWithPreviousIncorrect() {
        receiver.receive("678IDS".toByteArray(), "678IDS".toByteArray().size)
        receiver.receive("$434á438#".toByteArray(), "$434á438#".toByteArray().size)

        receiver.receive("$10#1#11#123".toByteArray(), "$10#1#11#123".toByteArray().size)
        receiver.receive("456789X".toByteArray(), "456789X".toByteArray().size)

        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)
    }

    @Test
    internal fun receiveMultipleValidMessages() {
        receiver.receive("$10#1#11#123".toByteArray(), "$10#1#11#123".toByteArray().size)
        receiver.receive("456789X".toByteArray(), "456789X".toByteArray().size)

        receiver.receive("$17#2#12#123".toByteArray(), "$17#2#12#123".toByteArray().size)
        receiver.receive("456789X".toByteArray(), "456789X".toByteArray().size)
        receiver.receive("456789X".toByteArray(), "456789X".toByteArray().size)

        assertEquals(2, received.size)

        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)

        assertEquals(17, received[1].length)
        assertEquals("123456789X456789X", received[1].content)
        assertEquals(2, received[1].type)
        assertEquals(12, received[1].id)
    }

    @Test
    internal fun receiveSingleInvalidMessage() {
        receiver.receive("$19#1#123".toByteArray(), "$19#1#123".toByteArray().size)
        receiver.receive("456789X".toByteArray(), "456789X".toByteArray().size)

        assertEquals(0, received.size)
    }

    @Test
    internal fun receiveMultipleInvalidMessage() {
        receiver.receive("1546f#1#123".toByteArray(), "1546f#1#123".toByteArray().size)
        receiver.receive("1546SDs\$dsdsd23".toByteArray(), 54)
        receiver.receive("1546f#1#12$3dsd".toByteArray(), 132)
        receiver.receive("1546f#d1465sad 46s5f1#123".toByteArray(), 12)

        assertEquals(0, received.size)
    }
}
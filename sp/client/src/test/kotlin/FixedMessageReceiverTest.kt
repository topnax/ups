import networking.reader.MessageReader
import networking.receiver.FixedMessageReceiver
import networking.receiver.Message
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test

class FixedMessageReceiverTest {

    private val received = mutableListOf<Message>()
    private lateinit var receiver: FixedMessageReceiver

    @BeforeEach
    internal fun setUp() {
        received.clear()
        receiver = FixedMessageReceiver(object : MessageReader {
            override fun read(message: Message) {
                received.add(message)
            }
        })
    }

    @Test
    internal fun splitMessagesTest() {
        val message = "$10#1#11#123456789X$9#1#11#123456789"
        val received = receiver.receive(message.toByteArray(), message.toByteArray().size)
        assertEquals(2, received.size)
        assertEquals("$10#1#11#123456789X", String(received[0]))
        assertEquals("$9#1#11#123456789", String(received[1]))
    }

    @Test
    internal fun splitMessagesTest2() {
        val message = """$120#111#0#{"tiles":[{"row":10,"column":6,"set":true,"highlighted":true,"type":0,"letter":{"value":""""
        val message2 = """ě","points":5,"PlayerID":1}}]}$33#701#10#{"content":"Placed successfully"}"""

        val received = receiver.receive(message.toByteArray(), message.toByteArray().size)
        val received2 = receiver.receive(message2.toByteArray(), message2.toByteArray().size)

        assertEquals(1, received.size)
        assertEquals(2, received2.size)
        assertEquals("""$120#111#0#{"tiles":[{"row":10,"column":6,"set":true,"highlighted":true,"type":0,"letter":{"value":"""", String(received[0]))
        assertEquals("""ě","points":5,"PlayerID":1}}]}""", String(received2[0]))
        assertEquals("""$33#701#10#{"content":"Placed successfully"}""", String(received2[1]))
        assertEquals("""{"tiles":[{"row":10,"column":6,"set":true,"highlighted":true,"type":0,"letter":{"value":"ě","points":5,"PlayerID":1}}]}""", this.received[0].content)
        assertEquals("""{"content":"Placed successfully"}""", this.received[1].content)
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
    internal fun simpleReceiveTest3() {
        val message = "$10#1#11#123456789X$19$9#1#11#123456"
        val messagePart2 = "987"
        val received = receiver.receive(message.toByteArray(), message.toByteArray().size)
        assertEquals("$10#1#11#123456789X", String(received[0]))
        assertEquals("$19", String(received[1]))
        assertEquals("$9#1#11#123456", String(received[2]))
        assertEquals(3, received.size)

        val received2 = receiver.receive(messagePart2.toByteArray(), messagePart2.toByteArray().size)
        assertEquals("987", String(received2[0]))

        assertEquals(2, this.received.size)
        assertEquals("123456789X", this.received[0].content)
        assertEquals("123456987", this.received[1].content)

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
    internal fun receiveSimpleOverriddenMessage() {
        val message = "$10#1#11#123$10#1#11#123456789X456789X"
        receiver.receive(message.toByteArray(), message.length)
        assertEquals(1, received.size)
        assertEquals(10, received[0].length)
        assertEquals("123456789X", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)
    }

    @Test
    internal fun receiveSimpleEscapedMessage() {
        val message = "$12#1#11#12345\\$6789X"
        receiver.receive(message.toByteArray(), message.length)
        assertEquals(1, received.size)
        assertEquals(12, received[0].length)
        assertEquals("12345\\\$6789X", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)
    }

    @Test
    internal fun receiveSimpleEscapedMessage2() {
        val message = "$12#1#11#12345\\\\$12#1#11#12345\\\\\\$678"
        receiver.receive(message.toByteArray(), message.length)
        assertEquals(1, received.size)
        assertEquals(12, received[0].length)
        assertEquals("12345\\\\\\$678", received[0].content)
        assertEquals(1, received[0].type)
        assertEquals(11, received[0].id)
    }

    @Test
    internal fun receiveSimpleEscapedMessage3() {
        val message = "$12#1#11#12345\\\\\\$678"
        receiver.receive(message.toByteArray(), message.length)
        assertEquals(1, received.size)
        assertEquals(12, received[0].length)
        assertEquals("12345\\\\\\$678", received[0].content)
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
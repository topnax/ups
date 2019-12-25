import networking.receiver.indexOfNth
import networking.receiver.isNextByteEscaped
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertFalse
import org.junit.jupiter.api.Test

class IsNextCharacterEscapedTest {

    @Test
    internal fun simpleNonEscapedTest() {
        assertFalse {
            "not escaped".toByteArray().toMutableList().isNextByteEscaped()
        }
    }

    @Test
    internal fun nonEscapedTest() {
        assertFalse {
            "not escaped\\\\".toByteArray().toMutableList().isNextByteEscaped()
        }
    }


    @Test
    internal fun nonEscapedTest2() {
        assertFalse {
            "not escaped\\\\\\\\".toByteArray().toMutableList().isNextByteEscaped()
        }
    }

    @Test
    internal fun simpleEscapedTest() {
        assert("escaped\\".toByteArray().toMutableList().isNextByteEscaped())
    }

    @Test
    internal fun simpleEscapedTest2() {
        assert("escaped\\\\\\".toByteArray().toMutableList().isNextByteEscaped())
    }

    @Test
    internal fun simpleEscapedTest3() {
        assert("escaped\\\\\\\\\\".toByteArray().toMutableList().isNextByteEscaped())
    }

}
import networking.receiver.indexOfNth
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Test

class IndexOfNthTest {

    @Test
    internal fun indexOfFirst() {
        assertEquals(3, "foobar".indexOfNth('b', 1))
    }

    @Test
    internal fun indexOfFirst2() {
        assertEquals(5, "foobar".indexOfNth('r', 1))
    }

    @Test
    internal fun indexOfSecond() {
        assertEquals(6, "foobarbaz".indexOfNth('b', 2))
    }

    @Test
    internal fun indexOfSecond2() {
        assertEquals(6, "foobarbabbbb".indexOfNth('b', 2))
    }

    @Test
    internal fun indexOfSecond3() {
        assertEquals(8, "foobarbabbbb".indexOfNth('b', 3))
    }

    @Test
    internal fun indexOfNotFound() {
        assertEquals(-1, "foobarbaz".indexOfNth('x', 3))
    }

    @Test
    internal fun indexOfNotFound2() {
        assertEquals(-1, "foobarbaz".indexOfNth('x', 0))
    }

    @Test
    internal fun indexOfNotFound3() {
        assertEquals(-1, "foobarbaz".indexOfNth('x', 1))
    }

}
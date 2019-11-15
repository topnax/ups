package networking.messages

import com.beust.klaxon.FieldRenamer
import com.beust.klaxon.Json
import com.beust.klaxon.Klaxon

abstract class ApplicationMessage(@Json(ignored = true) val type: Int) {

    companion object {
        private val renamer = object: FieldRenamer {
            override fun toJson(fieldName: String) = FieldRenamer.camelToUnderscores(fieldName)
            override fun fromJson(fieldName: String) = FieldRenamer.underscoreToCamel(fieldName)
        }
    }

    fun toJson(): String {
        return Klaxon().fieldRenamer(renamer).toJsonString(this)
    }
}

data class JoinLobbyMessage(val lobbyId: Int, val playerName: String, val foo: Foo) : ApplicationMessage(1)

data class Foo(val bar: Int, val baz: String)

fun main() {
    println(JoinLobbyMessage(12, "Topnax", Foo(42, "This is Foo")).toJson())
}
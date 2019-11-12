package hello

import java.io.*
import java.net.InetAddress
import java.net.Socket

class ChatUser(
    private val port: Int = 10000,
    private val hostname: String = "localhost",
    val clientName: String = "NoName"
) : Thread() {

    val socket: Socket

//    var input: ObjectInputStream = null
    val output: ObjectOutputStream
    val input: ObjectInputStream
    val stdin: BufferedReader

    init {
        println("Client opening a socket at $hostname using port $port")
        socket = Socket(InetAddress.getByName(hostname), port)
        println("Socket created @ ${socket.inetAddress} with port ${socket.port}")
        output = ObjectOutputStream(socket.getOutputStream())
        println("output gathered")
        input = ObjectInputStream(socket.getInputStream())
        println("output gathered")
        stdin = BufferedReader(InputStreamReader(System.`in`))
        println("init finished")

    }

    override fun run() {

        println("Inside run")
        var serverMessage: String? = null
        var stdinMessage: String? = null
        try {
            println("Writing to server")
            output.writeObject("Hello from client, my name is $clientName")
            output.writeObject("Hello from client, my name is $clientName")
            output.writeObject("Hello from client, my name is $clientName")

            while ({
                    stdinMessage = stdin.readLine() as String; stdinMessage
                }() != null || true) {
                println("inside")
                serverMessage.let {
                    println("from server: '$it'")
                }
                stdinMessage.let {
                    println("from stdin: '$it'")
                    output.writeObject(it)
                }
            }
        } catch (e: IOException) {
            e.printStackTrace()
        } finally {
            println("finished")
            finish()
        }
    }

    private fun finish() {
//        input.close()
        output.close()
    }
}
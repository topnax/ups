package hello

fun main(args: Array<String>) {
    if (args.size == 2 && args[0] == "s" && args[1].matches(Regex("-?\\d+"))) {
        println("starting server...")
        println(args[0])
        println(args[1])
        server(args[1].toInt())
    } else if (args.size == 3 && args[0] == "c" && args[1].matches(Regex("-?\\d+"))) {
        println("starting client...")
        println(args[0])
        println(args[1])
        println(args[2])
        SelectClient(args[1].toInt(), userName = args[2]).start()
    } else {
        println("Unknown arguments...")
    }
}

fun server(port: Int, hostname: String = "localhost") {
    SelectServer(port, hostname).start()
}

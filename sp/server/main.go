package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"syscall"
	"ups/sp/server/encoding"
	"ups/sp/server/game_server"
	"ups/sp/server/networking/server"
)

func main() {

	bytes, err := json.Marshal(encoding.JoinLobbyMessage{
		LobbyID:    1,
		ClientName: "Standa",
	})
	log.Infoln(string(bytes))

	//
	//messageReader := encoding.SimpleMessageReader{}
	//
	//messageReader.Receive(1, []byte("#17#1#{'name':'"))
	//
	//
	//messageReader.Receive(2, []byte("#16#2#{'name':'Pavel'}"))
	//
	//messageReader.Receive(1, []byte("Standa'}"))

	//
	////messageReader.Receive(1, []byte("ciao  girls"))
	//
	//
	//fufu := Foo{
	//	Name: "Stenly",
	//	Attr: Bar{
	//		Age: 21,
	//	},
	//}
	//
	//bubu := Bar{Age:12,}
	//
	//eat(fufu, 0)
	//
	//eat(bubu, 1)

	//sample := json.Unmarshal(rawMessage.)
	//
	//
	//fmt.Println(sample.Name)
	//fmt.Println(sample.Number)

	log.SetLevel(log.DebugLevel)
	//
	//log.SetOutput(os.Stdout)

	serverx, err := server.NewServer(syscall.SockaddrInet4{
		Addr: [4]byte{byte(127), byte(0), byte(0), byte(1)},
		Port: 10000,
	})

	if err != nil {
		log.Errorln(err)
		return
	}
	srdr := encoding.SimpleMessageReader{}

	srdr.SetResponseOutput(&serverx)

	kkmr := game_server.NewKrisKrosServer()

	jsreade := encoding.SimpleJsonReader{}
	jsreade.Init()
	jsreade.SetOutput(&kkmr)

	srdr.SetOutput(&jsreade)

	kkmr.SetMessageSender(&srdr)

	serverx.SetOnClientDisconnectedListener(&kkmr)

	serverx.Start(&srdr)

	//kkmr := encoding.KrisKrosServer{}
	//
	//jsreade := encoding.SimpleJsonReader{}
	//jsreade.Init()
	//jsreade.SetOutput(&kkmr)
	//
	//msg := encoding.SimpleMessage{
	//	ClientUID: 1,
	//	Length:    10,
	//	Type:      1,
	//	Content:   "{\"name\":\"Standa\", \"age\": 21}",
	//}
	//
	//jsreade.Read(msg)
	//
	//msg = encoding.SimpleMessage{
	//	ClientUID: 1,
	//	Length:    10,
	//	Type:      2,
	//	Content:   "{\"surname\":\"Král\", \"smr\": {\"name\":\"Standa\", \"age\": 21}}",
	//}
	//
	//logrus.Infoln("mymes len", msg.Length)
	//jsreade.Read(msg)

	//currentGame := game.Game{}
	//currentGame.AddPlayer("Pavel")
	//currentGame.AddPlayer("Tomáš")
	//currentGame.AddPlayer("Fanda")
	//
	//currentGame.Start()
	//
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      1,
	//	Column:   1,
	//	Letter:   "S",
	//})
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      2,
	//	Column:   1,
	//	Letter:   "E",
	//})
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      3,
	//	Column:   1,
	//	Letter:   "X",
	//})
	//currentGame.Print()
	//
	//currentGame.Next()
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      1,
	//	Column:   3,
	//	Letter:   "P",
	//})
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      2,
	//	Column:   3,
	//	Letter:   "E",
	//})
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      3,
	//	Column:   3,
	//	Letter:   "S",
	//})
	//
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      2,
	//	Column:   2,
	//	Letter:   "Y",
	//})
	//currentGame.Print()
	//fmt.Println(currentGame.AcceptTurn(currentGame.Players[0]))
	//fmt.Println(currentGame.AcceptTurn(currentGame.Players[1]))
	//fmt.Println(currentGame.AcceptTurn(currentGame.Players[2]))
	//currentGame.Next()
	//
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      2,
	//	Column:   4,
	//	Letter:   "S",
	//})
	//
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      4,
	//	Column:   1,
	//	Letter:   "Y",
	//})
	//
	//_ = currentGame.HandleSetAtEvent(game.SetLetterAtEvent{
	//	PlayerID: currentGame.CurrentPlayer.ID,
	//	Row:      5,
	//	Column:   1,
	//	Letter:   "H",
	//})
	//
	//fmt.Println(currentGame.AcceptTurn(currentGame.Players[0]))
	//fmt.Println(currentGame.AcceptTurn(currentGame.Players[1]))
	//fmt.Println(currentGame.AcceptTurn(currentGame.Players[2]))
	//currentGame.Print()

	////
	//desk := model.GetDesk()
	//
	//desk.Print()
	//desk.SetAt("A", 0, 1)
	//desk.SetAt("B", 2, 1)
	//desk.SetAt("B", 2, 1)
	//desk.SetAt("CH", 14, 14)
	//fmt.Println("Out of bounds")
	//desk.SetAt("K", 3, 3)
	//desk.SetAt("R", 3, 4)
	//desk.SetAt("I", 3, 5)
	//desk.SetAt("S", 3, 6)
	//
	//desk.SetAt("K", 2, 4)
	//
	//desk.SetAt("O", 4, 4)
	//desk.SetAt("S", 5, 4)
	//
	//desk.SetAt("Z", 14, 15)
	//desk.SetAt("Z", 15, 14)
	//desk.SetAt("CH", 14, 14)
	//desk.SetAt("131", 13, 13)
	//desk.SetAt("1", 12, 13)
	//desk.SetAt("Z", -1, 15)
	//desk.SetAt("Z", 15, -1)
	//desk.SetAt("Z", 15, -1)
	//
	//error := desk.SetWordAt(3, 3, 3, 6, 1)
	//error2 := desk.SetWordAt(2, 4, 5, 4, 1)
	//
	////desk.SetAt("B", 16,15)
	////desk.SetAt("B", 15,15)
	//desk.Print()
	//
	//if error != nil {
	//	fmt.Println(error)
	//}
	//
	//if error2 != nil {
	//	fmt.Println(error2)
	//}
}

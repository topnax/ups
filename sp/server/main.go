package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"syscall"
	"ups/sp/server/networking/server"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	serverx := server.Server{}
	err := serverx.Init(syscall.SockaddrInet4{
		Addr: [4]byte{byte(127), byte(0), byte(0), byte(1)},
		Port: 10000,
	})

	if err != nil {
		log.Errorln(err)
		return
	}

	serverx.Start()
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

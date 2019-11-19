package kris_kros_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ups/sp/server/game"
	"ups/sp/server/model"
	"ups/sp/server/rework/protocol/def"
	"ups/sp/server/rework/protocol/impl"
	"ups/sp/server/rework/protocol/messages"
)

type KrisKrosServer struct {
	count  int
	Router KrisKrosRouter
	sender def.ResponseSender
}

func failedToCast(message def.MessageHandler) def.Response {
	return impl.ErrorResponse(fmt.Sprintf("Failed to cast message of type %d", message.GetType()), impl.FailedToCast)
}

func NewKrisKrosServer(sender def.ResponseSender) KrisKrosServer {
	kks := KrisKrosServer{}
	kks.Router = newKrisKrosRouter(&kks)
	kks.sender = sender
	return kks
}

func (k KrisKrosServer) Read(message def.MessageHandler, clientUID int) def.Response {
	res := k.Router.route(message, clientUID)
	log.Infoln("Server returning: ", res.Content())
	return res
}

func (k *KrisKrosServer) OnPlayerJoined(message messages.PlayerJoinedMessage, clientUID int) def.Response {
	log.Infoln("On player joined:)", message.Name, k.count)
	k.count++
	k.sender.Send(impl.MessageResponse(messages.PlayerConfirmation{
		Lobby: model.Lobby{
			ID:    k.count,
			Owner: game.Player{Name: message.Name},
		},
	}, 101), clientUID)
	return impl.SuccessResponse("Successfully :)")
}

func (k *KrisKrosServer) onPlayerLeft(message def.MessageHandler) {
	println("opl")
}

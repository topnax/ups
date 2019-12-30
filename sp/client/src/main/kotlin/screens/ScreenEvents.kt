package screens

import networking.messages.GameStateRegenerationResponse
import tornadofx.EventBus
import tornadofx.FXEvent

class ServerRestartedEvent : FXEvent(EventBus.RunOn.BackgroundThread)

class ServerRestartedUnauthorizedEvent : FXEvent(EventBus.RunOn.BackgroundThread)

class DisconnectedEvent : FXEvent(EventBus.RunOn.BackgroundThread)

class ServerUnreachableEvent : FXEvent(EventBus.RunOn.BackgroundThread)

class GameRegeneratedEvent(val response: GameStateRegenerationResponse) : FXEvent(EventBus.RunOn.BackgroundThread)


class UserAuthenticatedEvent(val name: String) : FXEvent(EventBus.RunOn.BackgroundThread)



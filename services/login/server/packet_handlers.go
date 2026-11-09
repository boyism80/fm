package server

func (ls *LoginServer) registerPacketHandlers() {
	ls.Bind[Pong]()
	ls.Bind[Login]()
	ls.Bind[CharacterList]()
	ls.Bind[CheckName]()
	ls.Bind[CreateCharacter]()
	ls.Bind[DeleteCharacter]()
	ls.Bind[SelectCharacter]()
	ls.Bind[AuthSecondPassword]()
	ls.Bind[SelectCharacterWithSecondPassword]()
}

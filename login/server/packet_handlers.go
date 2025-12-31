package server

import (
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

// registerPacketHandlers registers all login server packet handlers
func (ls *LoginServer) registerPacketHandlers() {
	// Register the 7 login server packet handlers
	ls.server.RegisterPacketHandler(0x0A, ls.handlePong)
	ls.server.RegisterPacketHandler(0x01, ls.handleLogin)
	ls.server.RegisterPacketHandler(0x04, ls.handleCharacterList)
	ls.server.RegisterPacketHandler(0x07, ls.handleCheckName)
	ls.server.RegisterPacketHandler(0x08, ls.handleCreateCharacter)
	ls.server.RegisterPacketHandler(0x09, ls.handleDeleteCharacter)
	ls.server.RegisterPacketHandler(0x05, ls.handleSelectCharacter)

	log.Printf("Registered %d login server packet handlers", ls.server.GetPacketHandler().GetHandlerCount())
}

// handlePong processes pong packets from clients
func (ls *LoginServer) handlePong(ctx *core.ClientContext, data []byte) error {
	// Deserialize pong packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.Pong{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize pong packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Pong packet received from %s", ctx.Client.GetConnection().RemoteAddr())
	return nil
}

// handleLogin processes login authentication requests
func (ls *LoginServer) handleLogin(ctx *core.ClientContext, data []byte) error {
	// Deserialize login packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.Login{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize login packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Login packet received from %s - ID: %s, MAC: %s",
		ctx.Client.GetConnection().RemoteAddr(), request.ID, request.Mac)

	// Check if user is "cshyeon" (hardcoded for demo)
	if request.ID == "cshyeon" {
		// Send authentication success
		authResp := &response.Authenticate{
			AccountId:     2390,
			Gender:        0,
			Admin:         true,
			AccountName:   request.ID,
			IsChatBlocked: false,
			ChatBlockTime: 116445060000000000,
		}
		if err := ctx.SendFunc(authResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send authenticate response: %v", err)
			return err
		}

		// Send server list (10 servers for demo)
		for i := 0; i < 10; i++ {
			serverResp := &response.ServerList{
				ServerId:     uint8(i),
				ChannelSize:  5,
				WorldName:    "cshyeon",
				Flag:         0,
				EventMessage: "채승현이 간다.",
			}
			if err := ctx.SendFunc(serverResp, types.SEND_POLICY_ENCRYPT); err != nil {
				log.Printf("Failed to send server list response: %v", err)
				return err
			}
		}

		// Send end of server list
		endResp := &response.EndOfServerList{}
		if err := ctx.SendFunc(endResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send end of server list: %v", err)
			return err
		}
	} else {
		// Send login failed
		failedResp := &response.LoginFailed{Reason: response.LoginFailedReasonNoPopup}
		if err := ctx.SendFunc(failedResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send login failed response: %v", err)
			return err
		}

		// Send notice
		noticeResp := &response.Notice{
			Type:    constant.MSG_POPUP,
			Channel: 0,
			Message: "Hello",
			MegaEar: false,
		}
		if err := ctx.SendFunc(noticeResp, types.SEND_POLICY_ENCRYPT); err != nil {
			log.Printf("Failed to send notice: %v", err)
			return err
		}
	}

	return nil
}

// handleCharacterList processes character list requests
func (ls *LoginServer) handleCharacterList(ctx *core.ClientContext, data []byte) error {
	// Deserialize character list packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.CharacterList{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize character list packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Character list packet received from %s - Server: %d, Channel: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.Server, request.Channel)

	// Create dummy characters for demo (as DTO)
	characters := []dto.Character{
		{
			CharacterOverview: dto.CharacterOverview{
				ID:         1,
				Name:       "채승현",
				Gender:     0,
				SkinColor:  0,
				Face:       20100,
				Hair:       30000,
				Level:      255,
				Class:      0,
				Str:        12,
				Dex:        5,
				Int:        4,
				Luk:        4,
				Hp:         50,
				MaxHp:      50,
				Mp:         5,
				MaxMp:      5,
				SpawnPoint: 1,
			},
			CharacterLook: dto.CharacterLook{
				Gender:     0,
				SkinColor:  0,
				Face:       20100,
				Hair:       30000,
				Equipments: make(map[int8]uint32),
				Skins:      make(map[int8]uint32),
			},
		},
		{
			CharacterOverview: dto.CharacterOverview{
				ID:         2,
				Name:       "채진영",
				Gender:     0,
				SkinColor:  0,
				Face:       20100,
				Hair:       30000,
				Level:      255,
				Class:      0,
				Str:        12,
				Dex:        5,
				Int:        4,
				Luk:        4,
				Hp:         50,
				MaxHp:      50,
				Mp:         5,
				MaxMp:      5,
				SpawnPoint: 1,
			},
			CharacterLook: dto.CharacterLook{
				Gender:     0,
				SkinColor:  0,
				Face:       20100,
				Hair:       30000,
				Equipments: make(map[int8]uint32),
				Skins:      make(map[int8]uint32),
			},
		},
	}

	// Send character list response
	charListResp := &response.CharacterList{
		Characters: characters,
		SlotCount:  6,
	}
	if err := ctx.SendFunc(charListResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send character list response: %v", err)
		return err
	}

	return nil
}

// handleCheckName processes character name availability checks
func (ls *LoginServer) handleCheckName(ctx *core.ClientContext, data []byte) error {
	// Deserialize check name packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.CheckName{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize check name packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Check name packet received from %s - Name: %s",
		ctx.Client.GetConnection().RemoteAddr(), request.Name)

	// Check if name exists (hardcoded for demo)
	exists := request.Name == "채승현"

	// Send check name response
	checkResp := &response.CheckName{
		Name:   request.Name,
		Exists: exists,
	}
	if err := ctx.SendFunc(checkResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send check name response: %v", err)
		return err
	}

	return nil
}

// handleCreateCharacter processes character creation requests
func (ls *LoginServer) handleCreateCharacter(ctx *core.ClientContext, data []byte) error {
	// Deserialize create character packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.CreateCharacter{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize create character packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Create character packet received from %s - Name: %s, Face: %d, Hair: %d, Top: %d, Bottom: %d, Shoes: %d, Weapon: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.Name, request.Face, request.Hair, request.Top, request.Bottom, request.Shoes, request.Weapon)

	// Check if creation is successful (hardcoded for demo)
	success := request.Name != "채진영"

	// Create character response (as DTO)
	createResp := &response.CreateCharacter{
		Success: success,
		Character: &dto.Character{
			CharacterOverview: dto.CharacterOverview{
				ID:         1,
				Name:       request.Name,
				Gender:     0,
				SkinColor:  0,
				Face:       request.Face,
				Hair:       request.Hair,
				Level:      1,
				Class:      0,
				Str:        12,
				Dex:        5,
				Int:        4,
				Luk:        4,
				Hp:         50,
				MaxHp:      50,
				Mp:         5,
				MaxMp:      5,
				SpawnPoint: 3,
			},
			CharacterLook: dto.CharacterLook{
				Gender:     0,
				SkinColor:  0,
				Face:       request.Face,
				Hair:       request.Hair,
				Equipments: make(map[int8]uint32),
				Skins:      make(map[int8]uint32),
			},
		},
	}
	if err := ctx.SendFunc(createResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send create character response: %v", err)
		return err
	}

	return nil
}

// handleDeleteCharacter processes character deletion requests
func (ls *LoginServer) handleDeleteCharacter(ctx *core.ClientContext, data []byte) error {
	// Deserialize delete character packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.DeleteCharacter{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize delete character packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Delete character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.ID)

	// Send delete character response
	deleteResp := &response.DeleteCharacter{
		ID:      request.ID,
		Success: true,
	}
	if err := ctx.SendFunc(deleteResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send delete character response: %v", err)
		return err
	}

	return nil
}

// handleSelectCharacter processes character selection requests
func (ls *LoginServer) handleSelectCharacter(ctx *core.ClientContext, data []byte) error {
	// Deserialize select character packet
	reader := stream.NewStreamReader(&data, stream.LittleEndian)
	request := &request.SelectCharacter{}
	if err := request.Deserialize(reader); err != nil {
		log.Printf("Failed to deserialize select character packet from %s: %v", ctx.Client.GetConnection().RemoteAddr(), err)
		return err
	}

	log.Printf("Select character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), request.CharacterId)

	// Send transfer response to redirect to game server
	transferResp := &response.Transfer{
		IP:          ls.config.GameServerHost,
		Port:        uint16(ls.config.GameServerPort),
		CharacterId: request.CharacterId,
	}
	if err := ctx.SendFunc(transferResp, types.SEND_POLICY_ENCRYPT); err != nil {
		log.Printf("Failed to send transfer response: %v", err)
		return err
	}

	return nil
}

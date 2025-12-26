package net

// Opcode identifies the type of packet.
// Client->Server and Server->Client opcodes are separate namespaces.
//
// WoW uses 16-bit opcodes, grouped by category.
// We'll use a similar approach.
type Opcode uint16

// =============================================================================
// CLIENT -> SERVER OPCODES (CMSG = Client Message)
// =============================================================================

const (
	// Authentication (0x0000 - 0x00FF)
	CMSG_AUTH_SESSION    Opcode = 0x0001 // Initial auth with session token
	CMSG_PING            Opcode = 0x0002 // Keepalive
	CMSG_LOGOUT_REQUEST  Opcode = 0x0003
	CMSG_LOGOUT_CANCEL   Opcode = 0x0004

	// Character selection (0x0100 - 0x01FF)
	CMSG_CHAR_ENUM       Opcode = 0x0100 // Request character list
	CMSG_CHAR_CREATE     Opcode = 0x0101
	CMSG_CHAR_DELETE     Opcode = 0x0102
	CMSG_PLAYER_LOGIN    Opcode = 0x0103 // Enter world with character

	// Movement (0x0200 - 0x02FF)
	// These are sent frequently - optimize for size
	CMSG_MOVE_START_FORWARD  Opcode = 0x0200
	CMSG_MOVE_START_BACKWARD Opcode = 0x0201
	CMSG_MOVE_STOP           Opcode = 0x0202
	CMSG_MOVE_START_STRAFE_LEFT  Opcode = 0x0203
	CMSG_MOVE_START_STRAFE_RIGHT Opcode = 0x0204
	CMSG_MOVE_STOP_STRAFE    Opcode = 0x0205
	CMSG_MOVE_JUMP           Opcode = 0x0206
	CMSG_MOVE_FALL_LAND      Opcode = 0x0207
	CMSG_MOVE_HEARTBEAT      Opcode = 0x0208 // Periodic position update
	CMSG_MOVE_SET_FACING     Opcode = 0x0209

	// Combat (0x0300 - 0x03FF)
	CMSG_ATTACK_START    Opcode = 0x0300
	CMSG_ATTACK_STOP     Opcode = 0x0301
	CMSG_CAST_SPELL      Opcode = 0x0302
	CMSG_CANCEL_CAST     Opcode = 0x0303
	CMSG_CANCEL_AURA     Opcode = 0x0304

	// Interaction (0x0400 - 0x04FF)
	CMSG_SET_SELECTION   Opcode = 0x0400 // Target an entity
	CMSG_GOSSIP_HELLO    Opcode = 0x0401 // Talk to NPC
	CMSG_GOSSIP_SELECT   Opcode = 0x0402
	CMSG_QUEST_ACCEPT    Opcode = 0x0403
	CMSG_QUEST_COMPLETE  Opcode = 0x0404

	// Chat (0x0500 - 0x05FF)
	CMSG_CHAT_MESSAGE    Opcode = 0x0500
	CMSG_EMOTE           Opcode = 0x0501
)

// =============================================================================
// SERVER -> CLIENT OPCODES (SMSG = Server Message)
// =============================================================================

const (
	// Authentication responses (0x8000 - 0x80FF)
	SMSG_AUTH_RESPONSE   Opcode = 0x8001
	SMSG_PONG            Opcode = 0x8002
	SMSG_LOGOUT_RESPONSE Opcode = 0x8003
	SMSG_LOGOUT_COMPLETE Opcode = 0x8004

	// Character data (0x8100 - 0x81FF)
	SMSG_CHAR_ENUM       Opcode = 0x8100 // Character list
	SMSG_CHAR_CREATE     Opcode = 0x8101 // Create result
	SMSG_CHAR_DELETE     Opcode = 0x8102 // Delete result
	SMSG_LOGIN_VERIFY_WORLD Opcode = 0x8103 // Initial world position

	// Entity updates (0x8200 - 0x82FF)
	// These are the most frequent packets
	SMSG_UPDATE_OBJECT   Opcode = 0x8200 // Create/update entities
	SMSG_DESTROY_OBJECT  Opcode = 0x8201 // Entity left view
	SMSG_MOVE_UPDATE     Opcode = 0x8202 // Movement broadcast

	// Combat (0x8300 - 0x83FF)
	SMSG_ATTACK_START    Opcode = 0x8300
	SMSG_ATTACK_STOP     Opcode = 0x8301
	SMSG_ATTACKERSTATUS  Opcode = 0x8302
	SMSG_SPELL_START     Opcode = 0x8303
	SMSG_SPELL_GO        Opcode = 0x8304
	SMSG_SPELL_FAILURE   Opcode = 0x8305
	SMSG_DAMAGE_LOG      Opcode = 0x8306
	SMSG_HEAL_LOG        Opcode = 0x8307

	// Chat (0x8500 - 0x85FF)
	SMSG_CHAT_MESSAGE    Opcode = 0x8500
	SMSG_EMOTE           Opcode = 0x8501
	SMSG_SYSTEM_MESSAGE  Opcode = 0x8502
)

// TODO: Add opcode name lookup for debugging:
//
// var opcodeNames = map[Opcode]string{...}
//
// func (o Opcode) String() string

-- State machine (old/scripts/event/HenesysPQBonus.js): 헤네시스 PQ 보너스 (돼지의 마을)

local STAGE_MAP = 910010200
local EXIT_MAP = 910010300
local MOON_BUNNY_ID = 9300061
local DURATION_MS = 300000

local function clear(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		return { STAGE_MAP }
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(STAGE_MAP))
	end,

	on_mob_kill = function(sm, player, mobs)
		for _, mob in ipairs(mobs) do
			if mob ~= nil and mob:id() == MOON_BUNNY_ID then
				sm:message("월묘를 보호하지 못했습니다.", Msg.PinkText)
				clear(sm)
				return
			end
		end
	end,

	on_player_dead = function(sm, player)
	end,

	on_player_revive = function(sm, player)
	end,

	on_left_party = function(sm, player)
		clear(sm)
	end,

	on_disband_party = function(sm)
		clear(sm)
	end,

	on_scheduled_timeout = function(sm)
		clear(sm)
	end,

	on_all_monsters_dead = function(sm)
	end,
}

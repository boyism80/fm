-- State machine (old/scripts/event/HenesysPQ.js): 헤네시스 PQ (달맞이꽃 언덕)

local STAGE_MAP = 910010000
local EXIT_MAP = 910010300
local MOON_BUNNY_ID = 9300061
local RANKING_QUEST = 1200
local DURATION_MS = 600000

local function clear(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:set_property("stage", "0")
		group:set_property("clear", "0")
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		group:set_property("stage", "0")
		group:set_property("clear", "0")
		return { STAGE_MAP }
	end,

	on_prepare = function(sm)
		sm:map(STAGE_MAP):set_respawn(false)
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		local function on_arrive(player)
			player:try_party_quest(RANKING_QUEST)
		end
		player:map(sm:map(STAGE_MAP), { callback = on_arrive })
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

	on_clear = function(sm)
		clear(sm)
	end,

	on_finish = function(sm)
		local group = sm:group()
		group:set_property("state", "0")
		group:set_property("stage", "0")
		group:set_property("clear", "0")
	end,

	on_all_monsters_dead = function(sm)
	end,
}

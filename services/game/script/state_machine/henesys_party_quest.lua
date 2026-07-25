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
		group:declare_maps({ STAGE_MAP })
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_setup = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		group:set_property("stage", "0")
		group:set_property("clear", "0")
		local map = group:map(STAGE_MAP)
		map:reset()
		map:block_gen(false)
		sm:start_timer(DURATION_MS)
	end,

	on_player_entry = function(sm, player)
		player:map(STAGE_MAP)
		player:try_party_quest(RANKING_QUEST)
	end,

	on_mob_kill = function(sm, player, mobs)
		for _, mob in ipairs(mobs) do
			if mob ~= nil and mob:id() == MOON_BUNNY_ID then
				for _, p in ipairs(sm:players()) do
					if p ~= nil then
						p:notice("월묘를 보호하지 못했습니다.", Msg.PinkText)
					end
				end
				clear(sm)
				return
			end
		end
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == STAGE_MAP then
			return
		end
		sm:unregister(player)
	end,

	on_player_dead = function(sm, player)
	end,

	on_player_revive = function(sm, player)
	end,

	on_player_disconnected = function(sm, player)
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

	on_cancel_schedule = function(group)
	end
}

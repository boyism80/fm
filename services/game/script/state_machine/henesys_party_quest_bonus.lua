-- State machine (old/scripts/event/HenesysPQBonus.js): 헤네시스 PQ 보너스 (돼지의 마을)

local STAGE_MAP = 910010200
local EXIT_MAP = 910010300
local MOON_BUNNY_ID = 9300061
local DURATION_MS = 300000

function on_init(group)
	group:set_property("state", "0")
	group:declare_maps({ STAGE_MAP })
	group:declare_min_players(1)
	group:declare_exit_map(EXIT_MAP)
end

function on_setup(sm)
	local group = sm:group()
	group:set_property("state", "1")
	local map = group:map(STAGE_MAP)
	map:reset()
	sm:start_timer(DURATION_MS)
end

function on_player_entry(sm, player)
	player:map(STAGE_MAP)
end

function on_mob_kill(sm, player, mobs)
	for _, mob in ipairs(mobs) do
		if mob ~= nil and mob:id() == MOON_BUNNY_ID then
			for _, p in ipairs(sm:players()) do
				if p ~= nil then
					p:notice("월묘를 보호하지 못했습니다.", 5)
				end
			end
			on_clear(sm)
			return
		end
	end
end

function on_changed_map(sm, player, map_id)
	if map_id == STAGE_MAP then
		return
	end
	sm:unregister(player)
end

function on_player_dead(sm, player)
end

function on_player_revive(sm, player)
end

function on_player_disconnected(sm, player)
end

function on_left_party(sm, player)
	on_clear(sm)
end

function on_disband_party(sm)
	on_clear(sm)
end

function on_scheduled_timeout(sm)
	on_clear(sm)
end

function on_clear(sm)
	sm:finish(EXIT_MAP)
end

function on_finish(sm)
	sm:group():set_property("state", "0")
end

function on_all_monsters_dead(sm)
end

function on_cancel_schedule(group)
end

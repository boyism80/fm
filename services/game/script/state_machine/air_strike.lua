-- State machine (old/scripts/event/AirStrike.js): 바트의 방

local EXIT_MAP = 120000102
local QUEST_MAP = 912020000
local DURATION_MS = 300000

local function shuffle_reactors(map)
	local reactors = {}
	local positions = {}
	for _, reactor in pairs(map:reactors()) do
		local x, y = reactor:position()
		reactors[#reactors + 1] = reactor
		positions[#positions + 1] = { x, y }
	end
	for i = #positions, 2, -1 do
		local j = math.random(i)
		positions[i], positions[j] = positions[j], positions[i]
	end
	for i, reactor in ipairs(reactors) do
		reactor:position(positions[i][1], positions[i][2])
	end
end

local function finish(sm)
	local players = sm:players()
	sm:finish(0)
	for _, player in ipairs(players) do
		player:map(EXIT_MAP, 1)
	end
end

return {
	on_init = function(group)
		group:set_property("noEntry", "false")
		group:declare_min_players(1)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("noEntry", "true")
		local map = group:map(QUEST_MAP)
		if map ~= nil then
			map:reset()
			map:respawn(true)
			shuffle_reactors(map)
		end
		return { QUEST_MAP }

	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(QUEST_MAP)
	end,

	on_player_disconnected = function(sm, player)
		sm:finish(0)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("noEntry", "false")
	end
}

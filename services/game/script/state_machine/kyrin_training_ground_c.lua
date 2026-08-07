local TRAINING_MAP = 912010100
local EXIT_MAP = 120000101
local DURATION_MS = 240000

local function finish(sm, map_id)
	local players = sm:players()
	sm:finish(0)
	for _, player in ipairs(players) do
		player:map(map_id)
	end
end

return {
	on_init = function(group)
		group:set_property("started", "false")
		group:declare_maps({ TRAINING_MAP })
		group:declare_min_players(1)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("started", "true")
		local map = group:map(TRAINING_MAP)
		if map ~= nil then
			map:reset()
			map:respawn(true)
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(TRAINING_MAP)
		sm:notice("카이린의 공격으로부터 2분 이상 버티세요.", Msg.PinkText)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == TRAINING_MAP then
			return
		end
		sm:finish(0)
	end,

	on_player_disconnected = function(sm, player)
		sm:finish(0)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm, EXIT_MAP)
	end,

	on_finish = function(sm)
		sm:group():set_property("started", "false")
	end
}

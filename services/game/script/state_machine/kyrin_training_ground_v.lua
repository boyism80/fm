local TRAINING_MAP = 912010000
local EXIT_MAP = 912010200
local DURATION_MS = 180000

local function finish(sm)
	local players = sm:players()
	sm:finish(0)
	for _, player in ipairs(players) do
		player:map(EXIT_MAP)
	end
end

return {
	on_init = function(group)
		group:set_property("started", "false")
		group:min_players(1)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("started", "true")
		return { TRAINING_MAP }
	end,

	on_prepare = function(sm)
		sm:map(TRAINING_MAP):respawn({ include_one_time = true })
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(TRAINING_MAP))
		sm:message("카이린의 공격으로부터 2분 이상 버티세요.", Msg.PinkText)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("started", "false")
	end
}

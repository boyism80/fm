-- State machine (old/scripts/event/3rdjob.js): 3rd class trial

local DURATION_MS = 1200000

local TRIALS = {
	[105070001] = 108010300,
	[100040106] = 108010200,
	[105040305] = 108010100,
	[107000402] = 108010400,
	[105070200] = 108010500,
}

return {
	on_init = function(group)
		group:min_players(1)
		group:max_machines(nil)
	end,

	on_create = function(sm)
		local crack = sm:leader():map():template_id()
		local trial = TRIALS[crack]
		if trial == nil then
			sm:finish(0)
			return
		end
		sm:set_property("crack", tostring(crack))
		return { trial, trial + 1 }
	end,

	on_prepare = function(sm)
		local trial = TRIALS[tonumber(sm:get_property("crack"))]
		sm:map(trial + 1):respawn({ include_one_time = true })
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		local trial = TRIALS[tonumber(sm:get_property("crack"))]
		player:map(sm:map(trial), 0)
	end,

	on_scheduled_timeout = function(sm)
		sm:finish(tonumber(sm:get_property("crack")))
	end,
}

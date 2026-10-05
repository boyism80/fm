-- State machine (old/scripts/event/3rdjob.js): 3rd class trial

local DURATION_MS = 1200000
local ENTRANCE_PORTAL = "in00"

local TRIALS = {
	[105070001] = 108010300,
	[100040106] = 108010200,
	[105040305] = 108010100,
	[107000402] = 108010400,
	[105070200] = 108010500,
}

local function find_map(sm, template_id)
	for _, map in ipairs(sm:maps()) do
		if map:template_id() == template_id then
			return map
		end
	end
	return nil
end

return {
	on_init = function(group)
		group:declare_min_players(1)
	end,

	on_create = function(sm)
		local crack = sm:leader():map():template_id()
		local trial = TRIALS[crack]
		if trial == nil then
			sm:finish(0)
			return
		end
		local entrance, err = id2map(trial):create_instance()
		if entrance == nil then
			log("third_class_trial create_instance:", err)
			sm:finish(0)
			return
		end
		local arena, arena_err = id2map(trial + 1):create_instance()
		if arena == nil then
			log("third_class_trial create_instance:", arena_err)
			entrance:destroy()
			sm:finish(0)
			return
		end
		arena:respawn({ include_one_time = true })
		entrance:portal(ENTRANCE_PORTAL):script("instance_link")
		sm:set_property("crack", tostring(crack))
		return { entrance, arena }
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		local trial = TRIALS[tonumber(sm:get_property("crack"))]
		player:map(find_map(sm, trial), 0)
	end,

	on_scheduled_timeout = function(sm)
		sm:finish(tonumber(sm:get_property("crack")))
	end,

	on_finish = function(sm)
		for _, map in ipairs(sm:maps()) do
			map:destroy()
		end
	end,
}

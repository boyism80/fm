-- State machine (old/scripts/event/TamePig.js): 켄타의 사육실

local QUEST_MAP = 923010000
local EXIT_MAP = 923010100
local TAME_PIG_ID = 9300102
local DURATION_MS = 300000

local function finish(sm, message)
	local players = sm:players()
	sm:notice(message, Msg.PinkText)
	sm:finish(0)
	for _, player in ipairs(players) do
		player:map(EXIT_MAP)
	end
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps({ QUEST_MAP })
		group:declare_min_players(1)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		local map = group:map(QUEST_MAP)
		if map ~= nil then
			map:reset()
			map:respawn(true)
			for _, mob in pairs(map:mobs(TAME_PIG_ID)) do
				map:remove_mob(mob:oid(), MobDieAnimation.FadeOut)
			end
			map:spawn_mob(TAME_PIG_ID, -26, 335)
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(QUEST_MAP)
		sm:notice("멧돼지를 외계인으로부터 보호하고, 페로몬과 연구 보고서를 회수하세요!", Msg.PinkText)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == QUEST_MAP then
			return
		end
		sm:finish(0)
	end,

	on_mob_die = function(sm, mob)
		if mob ~= nil and mob:id() == TAME_PIG_ID then
			finish(sm, "멧돼지를 보호하는 데 실패하였습니다.")
		end
	end,

	on_player_disconnected = function(sm, player)
		sm:finish(0)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm, "제한시간이 다 되어 실패하였습니다.")
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end
}

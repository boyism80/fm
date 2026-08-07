local pq = require("script/lib/party_quest")

local ALTAIR_FRAGMENT = 4001198
local RANKING_QUEST = 1206
local EXIT_MAP = 930000800

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local empty = pq.mob_count(map) == 0
		local has_golem = pq.mob_count(map, 9300183) > 0
		local reactor = map:reactor_by_name("")
		local reactor_ok = reactor == nil or reactor:state() == 1
		if (empty or has_golem) and reactor_ok then
			local code = me:exchange({}, {
				item = { [ALTAIR_FRAGMENT] = 1 },
				exp = 52000,
			})
			if code == ExchangeResult.LackCapacity then
				me:notice("포이즌 골렘을 제거하고, 인벤토리 공간을 비워주세요.", Msg.PinkText)
				return
			end
			if code ~= ExchangeResult.OK then
				me:notice("포이즌 골렘을 제거하고, 인벤토리 공간을 비워주세요.", Msg.PinkText)
				return
			end
			me:end_party_quest(RANKING_QUEST)
			me:map(EXIT_MAP)
			return
		end
		me:notice("포이즌 골렘을 제거하고, 인벤토리 공간을 비워주세요.", Msg.PinkText)
	end
}

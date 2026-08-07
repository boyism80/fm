local pq = require("script/lib/party_quest")

local REWARD_ITEM = 4031448
local EXIT_MAP = 220070400

return {
	on_enter = function(me)
		if not pq.has_item(me, REWARD_ITEM) then
			local code = me:exchange({}, { item = { [REWARD_ITEM] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:notice("인벤토리 공간이 부족합니다.", Msg.PinkText)
				return
			end
			if code ~= ExchangeResult.OK then
				return
			end
		end
		local sm = me:state_machine()
		if sm ~= nil then
			sm:finish(0)
		end
		me:map(EXIT_MAP)
	end
}

local pq = require("script/lib/party_quest")

local GROUP_NAME = "phoenix_egg_firehawk"
local EGG_ID = 4001113

local function can_enter(me)
	local q6240 = me:quest(6240)
	local q6241 = me:quest(6241)
	if q6240:started() or q6241:started() then
		return true
	end
	return q6240:completed() and not q6241:started() and not q6241:completed()
end

return {
	on_enter = function(me)
		if not can_enter(me) then
			me:notice("알 수 없는 힘으로 봉인되어 있습니다.", Msg.PinkText)
			return
		end
		if pq.has_item(me, EGG_ID) then
			me:notice("이미 피닉스의 알을 갖고 있어 입장할 수 없습니다.", Msg.PinkText)
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:notice("지금은 시련의 동굴에 들어갈 수 없습니다.", Msg.PinkText)
			return
		end
		if group:get_property("started") == "true" then
			me:notice("이미 다른 누군가가 퀘스트에 도전 중입니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:notice("지금은 시련의 동굴에 들어갈 수 없습니다.", Msg.PinkText)
			if err ~= nil then
				log("phoenix_egg_firehawk start_solo:", err)
			end
		end
	end
}

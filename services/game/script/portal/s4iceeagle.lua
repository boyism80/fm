local pq = require("script/lib/party_quest")

local EGG_ID = 4001114
local TARGET_MAP = 921100210

local function can_enter(me)
	local q6242 = me:quest(6242)
	local q6243 = me:quest(6243)
	if q6242:started() or q6243:started() then
		return true
	end
	return q6242:completed() and not q6243:started() and not q6243:completed()
end

return {
	on_enter = function(me)
		if not can_enter(me) then
			me:notice("알 수 없는 힘으로 봉인되어 있습니다.", Msg.PinkText)
			return
		end
		if pq.has_item(me, EGG_ID) then
			me:notice("이미 프리져의 알을 갖고 있어 입장할 수 없습니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(TARGET_MAP)
	end
}

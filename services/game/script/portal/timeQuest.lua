local pq = require("script/lib/party_quest")

local DEEP_TEMPLE = 270040000
local FORGOTTEN_KEY = 4032002

local GATES = {
	[270010100] = { quest = 3501, pass = 270010110, back = 270010000 },
	[270010200] = { quest = 3502, pass = 270010210, back = 270010110 },
	[270010300] = { quest = 3503, pass = 270010310, back = 270010210 },
	[270010400] = { quest = 3504, pass = 270010410, back = 270010310 },
	[270010500] = { quest = 3507, pass = 270020000, back = 270010410 },
	[270020100] = { quest = 3508, pass = 270020110, back = 270020000 },
	[270020200] = { quest = 3509, pass = 270020210, back = 270020110 },
	[270020300] = { quest = 3510, pass = 270020310, back = 270020210 },
	[270020400] = { quest = 3511, pass = 270020410, back = 270020310 },
	[270020500] = { quest = 3514, pass = 270030000, back = 270020410 },
	[270030100] = { quest = 3515, pass = 270030110, back = 270030000 },
	[270030200] = { quest = 3516, pass = 270030210, back = 270030110 },
	[270030300] = { quest = 3517, pass = 270030310, back = 270030210 },
	[270030400] = { quest = 3518, pass = 270030410, back = 270030310 },
	[270030500] = { quest = 3521, pass = 270040000, back = 270030410 },
}

return {
	on_enter = function(me)
		local map_id = me:map():wz():id()
		if map_id == DEEP_TEMPLE then
			if pq.has_item(me, FORGOTTEN_KEY) == false then
				me:message("무언가가 막고 있는 듯 더 이상 들어갈 수 없다.")
				return
			end
			me:rmitem(FORGOTTEN_KEY, 1)
			me:message("신전 깊은 곳으로 이동합니다.")
			me:play_portal_sound()
			me:map(270040100, "out00")
			return
		end

		local gate = GATES[map_id]
		if gate == nil then
			return
		end
		me:play_portal_sound()
		if me:quest(gate.quest):completed() then
			me:map(gate.pass, "out00")
		else
			me:message("허가받지 않은 사람은 신전의 흐름을 역류할 수 없어 이전 장소로 되돌아 갑니다.")
			me:map(gate.back, 0)
		end
	end
}

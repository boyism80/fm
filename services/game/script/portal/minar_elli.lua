local pq = require("script/lib/party_quest")

local SEED = 4031346

return {
	on_enter = function(me)
		if pq.has_item(me, SEED) == false then
			me:message("이 포탈을 사용하려면 마법의 씨앗이 필요합니다.")
			return
		end
		me:rmitem(SEED, 1)
		me:message("마법의 씨앗의 힘으로 어딘가로 이동됩니다..")
		me:play_portal_sound()
		if me:map():wz():id() == 240010100 then
			me:map(101010000, "minar00")
		else
			me:map(240010100, "elli00")
		end
	end
}

local pq = require("script/lib/party_quest")

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		local map_id = wz.id
		if map_id == 930000000 then
			me:notice("엘린의 변신 마법이 몸 안으로 스며든다.", Msg.PinkText)
			me:play_portal_sound()
			me:map(930000010)
		elseif map_id == 930000010 then
			me:play_portal_sound()
			me:map(930000100)
		elseif map_id == 930000100 then
			if pq.mob_count(map) == 0 then
				me:play_portal_sound()
				me:map(930000200)
			else
				me:notice("모든 몬스터를 없애기 전에는 이동할 수 없습니다.", Msg.PinkText)
			end
		elseif map_id == 930000200 then
			local spine = map:reactor_by_name("spine")
			if spine ~= nil and spine:state() < 4 then
				me:notice("가시 덤불이 길을 막고 있습니다.", Msg.PinkText)
			else
				me:play_portal_sound()
				me:map(930000300)
			end
		end
	end
}

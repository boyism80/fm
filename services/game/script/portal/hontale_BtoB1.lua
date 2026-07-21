-- Portal (old/scripts/portal/hontale_BtoB1.js): 생명의동굴 미로방 입구

local pq = require("script/lib/party_quest")

function on_enter(me)
	if not pq.has_item(me, 4001087, 1) then
		me:notice("미로방에 들어가는데 필요한 열쇠가 없습니다.", Msg.PinkText)
		return
	end
	pq.remove_all(4001087, me)
	me:notice("첫 번째 미로방의 수정의 힘에 의해 어딘가로 이동됩니다.", Msg.PinkText)
	me:play_portal_sound()
	me:map(240050101)
end

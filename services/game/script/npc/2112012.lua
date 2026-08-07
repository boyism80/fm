-- NPC name (String.wz/Npc.img.xml): 유레테

local pq = require("script/lib/party_quest")
local rj = require("script/lib/romeo_juliet")

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		local sm = me:state_machine()
		if map_id == 926100500 or map_id == 926110500 then
			if sm == nil then
				return
			end
			if sm:get_property("persuade_urete") == "1" then
				me:dialog(npc, "이럴수가.. 내가 그동안 쌓아온 모든 것을 잃고 말았어...")
				me:dialog(npc, "그렇게 말해주니 무척 고마워. 앞으로는 마가티아에 협력하도록 하겠어.")
				map:message("유레테는 마음을 고쳐먹고 마가티아에 협력하겠다고 말한다.")
				pq.party_exp(sm, 10500)
				map:message("유레테를 구하는데 성공하여 추가 경험치가 지급됩니다.")
			end
			pq.party_warp(sm, map_id + 100)
			return
		end
		if map_id == 926100600 or map_id == 926110600 then
			rj.exchange_marbles(me, npc)
		end
	end
}

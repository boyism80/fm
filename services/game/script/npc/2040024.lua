-- NPC name (String.wz/Npc.img.xml): 첫번째 에오스의 돌

local SCROLL = 4001020

return {
	on_click = function(me, npc)
		local count = 0
		for _, it in pairs(me:item(SCROLL)) do
			count = count + it:count()
		end
		if count < 1 then
			me:dialog(npc, "첫번째 에오스의 돌이다. 하지만 #b#t4001020##k가 없어 활성화 시킬 수 없을 것 같다.")
			return
		end
		if not me:dialog_yes_no(npc, "#b#t4001020##k를 사용하여 #b#p2040025##k이 있는 71층으로 이동하시겠습니까?") then
			return
		end
		local code = me:exchange({ item = { [SCROLL] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			return
		end
		me:map(221022900, 3)
	end
}

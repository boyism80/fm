-- NPC name (String.wz/Npc.img.xml): 엘나스 마법석

local STONE = 4001019

return {
	on_click = function(me, npc)
		local count = 0
		for _, it in pairs(me:item(STONE)) do
			count = count + it:count()
		end
		if count < 1 then
			me:dialog(npc, "오르비스로 이동할 수 있을 것 같은 마법석이다. 하지만 #b#t4001019##k가 없으면 활성화 시킬 수 없을 것 같다.")
			return
		end
		if not me:dialog_yes_no(npc, "#b#t4001019##k를 사용하여 #b#m200080200##k으로 이동하시겠습니까?") then
			return
		end
		local code = me:exchange({ item = { [STONE] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			return
		end
		me:map(200080200)
	end
}

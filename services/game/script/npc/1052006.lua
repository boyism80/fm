-- NPC name (String.wz/Npc.img.xml): 웅이

local MESO = { 500, 1200, 2000 }
local TICKETS = { 4031036, 4031037, 4031038 }

return {
	on_click = function(me, npc)
		local level = me:level()
		if level <= 19 then
			me:dialog(npc, "안녕하세요~ 커닝시티 지하철 입니다. 아직 공사중인 구간이 많아 위험하답니다~ 흐음.. 아직 공사중 구간에 들어가시긴 너무 약해보이시는데요~", false, true)
			return
		end

		local options = { "#b공사장 B1#k" }
		if level > 29 then
			options[2] = "#b공사장 B2#k"
		end
		if level > 39 then
			options[3] = "#b공사장 B3#k"
		end
		local sel = me:dialog_list(npc, "안녕하세요~ 커닝시티 지하철 입니다. 아직 공사중인 구간이 많아 위험하답니다~ 어떤 구간의 입장권을 구매하시고 싶으세요?\r\n#b", options)
		if sel == nil then
			return
		end

		if not me:dialog_yes_no(npc, "음.. 정말 #b공사장 B" .. sel .. "#k 으로 가시는 입장권을 구매하시겠어요? 가격은 " .. MESO[sel] .. " 메소#k 랍니다.") then
			me:dialog(npc, "그런가요? 아직 위험한 구간이 많으니 신중하게 생각하시고 다시 말을 걸어주세요.", false, true)
			return
		end

		local code = me:exchange({ meso = MESO[sel] }, { item = { [TICKETS[sel]] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "요금이 부족하신것 같은데요. 혹은 인벤토리 공간이 부족하신 것 같은데요? 다시 한번 확인해 보세요.", false, true)
			return
		end
		me:dialog(npc, "다음에 또 이용해 주세요~", false, true)
	end
}

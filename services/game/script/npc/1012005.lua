-- NPC name (String.wz/Npc.img.xml): 클로이

local fireworks = require("script/lib/fireworks")
local pq = require("script/lib/party_quest")

local KEG = "keg"

return {
	on_click = function(me, npc)
		-- The legacy script gave the powder keg donation only to GMs; everyone else gets the pet master's greeting.
		if not pq.is_gm(me) then
			me:dialog(npc, "펫에 대해 궁금한 점이 있는가..?")
			return
		end

		local sel = me:dialog_list(npc, "폭죽을 만들려면 #b화약통#k이 필요해요. 모두가 모은 화약통이 가득 차면 폭죽을 쏠 수 있어요!", {
			"화약통을 가져왔어요.",
			"화약통이 얼마나 모였는지 보여주세요.",
		})
		if sel == nil then
			return
		end
		if sel == 2 then
			me:dialog(npc, "화약통 모집 현황: " .. fireworks.percentage(KEG) .. "%\r\n전부 모이면 폭죽을 쏠 수 있어요...")
			return
		end

		local count = me:dialog_number(npc, "화약통을 가져오셨나요? 가지고 계신 #b화약통#k을 모두 주세요. 멋진 폭죽을 만들게요. 몇 개를 주시겠어요?", 1, 1, 10000)
		if count == nil then
			return
		end
		if not fireworks.give(me, KEG, count) then
			me:dialog(npc, "화약통이 부족해요. 다시 확인해 주세요.")
			return
		end
		me:dialog(npc, "화약통을 구하시면 또 가져다 주세요.")
	end
}

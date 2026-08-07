-- NPC name (String.wz/Npc.img.xml): 카이린

local pq = require("script/lib/party_quest")

local EXIT_MAP = 120000101

return {
	on_click = function(me, npc)
		local map = me:map()
		local wz = map ~= nil and map:wz() or nil
		if wz == nil then
			return
		end
		local item_id = 4031857
		if wz.id >= 108000502 then
			item_id = 4031856
		end
		local message
		if pq.has_item(me, item_id, 15) then
			message = "호오, 벌써 #b#t" .. item_id .. "##k 15개를 다 모은 거야? 대단한걸? 조금은 다시 봤어. 이곳에서 내보내 줄까?"
		else
			message = "아직 #b#t" .. item_id .. "##k 15개를 다 모으지 못한 모양이지? 그런데 지금 바로 나가고 싶은 거야?"
		end
		if not me:dialog_yes_no(npc, message) then
			me:dialog(npc, "나가고 싶다면 언제든지 내게 말을 걸어.")
			return
		end
		local sm = me:state_machine()
		if sm ~= nil then
			sm:finish(0)
		end
		me:map(EXIT_MAP)
	end
}

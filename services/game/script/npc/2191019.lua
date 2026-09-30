-- NPC name (String.wz/Npc.img.xml): 자이로드롭

local COUNTER = 22000
local QUEST = 31304
local FUR = 4032653

return {
	on_click = function(me, npc)
		local counter = me:quest(COUNTER)
		if counter:started() == false then
			counter:start("0")
		end
		local count = tonumber(counter:record()) or 0
		if count >= 5 or me:quest(QUEST):started() == false then
			me:dialog(npc, "자이로 드롭이다.")
			return
		end
		if me:empty_slots(InventoryType.Etc) < 1 then
			me:dialog(npc, "기타창 한칸을 비워주세요")
			return
		end
		local sel = me:dialog_list(npc, "끼익 덜컹..", {
			"자이로드롭 꼭대기에 수상한 흰 털이 끼어있는 것을 발견했습니다. 빼낼까요?",
		})
		if sel == nil then
			return
		end
		counter:record(tostring(count + 1))
		me:exchange(nil, { item = { [FUR] = 1 } })
		me:dialog(npc, "자이로드롭이 작동되었습니다. 알시오네에게 보고하세요.")
	end
}

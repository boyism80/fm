-- NPC name (String.wz/Npc.img.xml): 카슨

local QUEST = 3310
local RESEARCH = 4031709

return {
	on_click = function(me, npc)
		if me:quest(QUEST):started() == false then
			me:dialog(npc, "제뉴미스트.. 진리를 탐구하는 생명연금의 자랑스러운 학파라네.")
			return
		end
		if next(me:item(RESEARCH)) ~= nil then
			me:dialog(npc, "이미 #b#t4031709##k를 구해 온 것 같구만. 그렇다면 이 안에 다시 들어갈 필요는 없어보이네.")
			return
		end
		if me:dialog_accept(npc, "흐음, 지금 바로 폐쇄된 연구실로 가보겠는가?") == false then
			me:dialog(npc, "아직 준비가 덜 된 모양인가? 준비가 되는 대로 나를 찾아오게.")
			return
		end
		local sm, err = state_machine("jenumist_homu"):start_solo(me)
		if sm == nil then
			log("jenumist_homu start_solo:", err)
			me:dialog(npc, "흐음, 이미 이 안에 다른 누군가가 퀘스트에 도전하고 있네. 나중에 다시 시도해보게.")
		end
	end
}

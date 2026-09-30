-- NPC name (String.wz/Npc.img.xml): 도적 전직교관

local DARK_MARBLE = 4031013
local PROOF = 4031012

return {
	on_click = function(me, npc)
		local marbles = 0
		for _, it in pairs(me:item(DARK_MARBLE)) do
			marbles = marbles + it:count()
		end
		if marbles >= 30 then
			if not me:dialog(npc, "호오, 대단하군. 이렇게나 많은 검은 구슬을 모아올 줄이야. 좋네. 자네를 인정해주겠네. 이것을 갖고 다크로드님께 돌아가 보게나.", false, true) then
				return
			end
			if me:exchange({ item = { [DARK_MARBLE] = marbles } }, { item = { [PROOF] = 1 } }) ~= ExchangeResult.OK then
				me:dialog(npc, "자네..인벤토리 공간이 부족한건 아닌가? 기타 탭의 공간을 한칸 이상 비운 후 다시 나에게 말을 걸게나.")
				return
			end
		else
			if not me:dialog_yes_no(npc, "흐음. 아직 검은 구슬 30개는 모으지 못한 것 같군. 아직 검은구슬은 다 모으지 못했지만 바깥으로 나가고 싶은가? 지금 포기하면 처음부터 다시 시작해야한다네.") then
				me:dialog(npc, "그렇지. 아직 포기하기엔 이르다네. 조금만 더 노력해 보게")
				return
			end
		end
		local left = 0
		for _, it in pairs(me:item(DARK_MARBLE)) do
			left = left + it:count()
		end
		if left > 0 then
			me:rmitem(DARK_MARBLE, left)
		end
		me:map(102040000)
	end
}

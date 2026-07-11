local quest_id = 3305

local function has_equipped(me, item_id)
	local cape = me:equipped(EquipmentPart.Cape)
	if cape ~= nil and cape:wz():id() == item_id then
		return true
	end
	return false
end

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if has_equipped(me, 1102135) then
		me:dialog(npc, "제뉴미스트의 망토가 이미 있는걸 보니 새로운 망토가 필요하지 않겠군.", false, false)
		return
	end

	local q3347 = me:quest(3347)
	if q3347 ~= nil and q3347:completed() then
		me:dialog(npc, "제뉴미스트로써 누릴 수 있는 혜택을 이미 다 받은 것 같은데... 유감이지만 더 이상 제뉴미스트의 망토를 재지급해 줄 수 없네.", false, false)
		return
	end

	if not me:dialog_accept(npc, "으흠? 자네는 얼마 전에 제뉴미스트 소속이 된 연금술사 아닌가? 그런데 제뉴미스트의 망토는 어디에... 이런. 망토를 잃어버린 모양이군. 망토는 자랑스러운 제뉴미스트 학파임을 증명하는 소중한 것. 그런데 그런 것을 잃어버리다니... 정신이 있는 겐가? 매우 불쾌하군. 자네 같은 사람에게는 다시 망토를 주고 싶지 않지만... 어쩔 수 없지. 망토를 다시 받겠는가?") then
		me:dialog(npc, "싫다면 하는 수 없지. 자네는 더 이상 제뉴미스트 소속이 아닐세.", false, false)
		return
	end

	me:dialog(npc, "한번 더 기회를 주도록 하지. 제뉴미스트의 망토를 만들고 싶다면, #b동물의 가죽 10개#k와 #b에메랄드 #k5개, 연성하는 사람을 위한 수고비 #b10000메소#k를 가져오게.", false, true)
	q:start(npc, true)
end

local quest_id = 3306

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
		me:dialog(npc, "알카드노의 망토를 이미 가지고 있는걸 보니 새로운 망토가 필요하지 않겠군.", false, false)
		return
	end

	local q3347 = me:quest(3347)
	if q3347 ~= nil and q3347:completed() then
		me:dialog(npc, "알카드노로써 누릴 수 있는 혜택을 이미 다 받지 않았나. 유감이지만 더 이상 알카드노의 망토를 재지급해 줄 수 없네.", false, false)
		return
	end

	if not me:dialog_accept(npc, "본 기억이 있는 얼굴인 걸 보니 자네는 분명 알카드노 소속의 연금술사인 모양인데... 왜 알카드노의 망토를 하고 있지 않은 건가? 알카드노 소속이라면 누구나 하고 있어야 한다고 말했는데... 뭐? 망토를 잃어버렸다고...? 더 이상 알카드노가 되고 싶지 않다는 말인가? 흐음... 그건 아닌 모양이군. 그럼 망토를 다시 지급할 테니 망토를 받을 텐가?") then
		me:dialog(npc, "싫다면 할 수 없지. 망토가 없다고 알카드노 소속이 아니게 되는 것은 아니지만... 누구도 자네를 알카드노로 생각하지 않을 것이니 명심하게.", false, false)
		return
	end

	me:dialog(npc, "알카드노에 처음 들어올 때가 망토를 그냥 지급했지만, 재지급할 때는 직접 재료를 모아와야 하네. 재료는 #b#t4000021# 10개#k와 #b#t4021006# 5개#k, 연성하는 사람에게 지급할 #b10000메소#k라네. 그럼 잊지 말고 가져오게.", false, true)
	q:start(npc, true)
end

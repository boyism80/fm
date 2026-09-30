-- NPC name (String.wz/Npc.img.xml): 샹크스

local LETTER = 4031801
local FARE = 150

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "이 배를 타면 더 넓은 대륙으로 내려갈 수 있지. 150메소에 #b빅토리아 아일랜드#k로 데려다 줄게. 대신 한 번 여길 떠나면 다시는 돌아올 수 없어. 어때? 빅토리아 아일랜드로 가고 싶어? 가고 싶다면 지금 바로 데려다 줄 수 있지.") then
			me:dialog(npc, "흠.. 아직 이 곳에서 할 일이 남았나 보지?")
			return
		end
		local letter = item_count(me, LETTER) >= 1
		if letter then
			if me:dialog(npc, "좋아. 그럼 어서 150메소를 줘... 응? 그건 암허스트의 장로 루카스님의 추천서잖아? 뭐야, 이런 게 있었으면 진작 말을 했어야지. 루카스님이 추천할 정도로 재능 있는 모험가에게 돈을 받을 정도로 이 샹크스가 매정하지는 않다고!", false, true) == false then
				return
			end
			if me:dialog(npc, "추천서를 가지고 있으니, 특별히 요금은 면제해 줄게. 이제 빅토리아 아일랜드로 출발한다! 흔들릴지도 모르니 꽉 잡아!", true, true) == false then
				return
			end
			if me:exchange({ item = { [LETTER] = 1 } }, nil) ~= ExchangeResult.OK then
				return
			end
			me:map(104000000)
			return
		end
		if me:dialog(npc, "여긴 지루해졌지? 그럼 일단 150 메소부터 받고...", false, true) == false then
			return
		end
		if me:meso() < FARE then
			me:dialog(npc, "뭐야? 돈도 없으면서 가겠다고 한거야? 이상한 녀석이로군!")
			return
		end
		if me:dialog(npc, "오케이! 150메소도 받았겠다. 지금 바로 출발한다, 꽉 잡아!", false, true) == false then
			return
		end
		if me:exchange({ meso = FARE }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "뭐야? 돈도 없으면서 가겠다고 한거야? 이상한 녀석이로군!")
			return
		end
		me:map(104000000)
	end
}

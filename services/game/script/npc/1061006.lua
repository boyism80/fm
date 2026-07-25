-- NPC name (String.wz/Npc.img.xml): 이상한 모양의 석상

return {
	on_click = function(me, npc)
		local function item_count(item_id)
			local count = 0
			for _, it in pairs(me:item(item_id)) do
				count = count + it:count()
			end
			return count
		end
		if item_count(4031025) >= 10 then
			me:dialog(npc, "이미 #t4031025#를 가지고 있는 것 같습니다.")
			return
		end
		if item_count(4031028) >= 30 then
			me:dialog(npc, "이미 #t4031028#를 가지고 있는 것 같습니다.")
			return
		end
		if item_count(4031026) >= 20 then
			me:dialog(npc, "이미 #t4031026#를 가지고 있는 것 같습니다.")
			return
		end
		local q2052 = me:quest(2052)
		local q2053 = me:quest(2053)
		local q2054 = me:quest(2054)
		local active = false
		if q2052 ~= nil and (q2052:started() or q2052:completed()) then
			active = true
		elseif q2053 ~= nil and (q2053:started() or q2053:completed()) then
			active = true
		elseif q2054 ~= nil and (q2054:started() or q2054:completed()) then
			active = true
		end
		if active then
			if not me:dialog_yes_no(npc, "석상에 손을 대자 어디론가 빨려드는듯한 느낌이 듭니다. 이대로 이동하시겠습니까?") then
				me:dialog(npc, "석상에서 손을 떼자 아무일도 없던 것처럼 원래대로 돌아왔습니다.")
				return
			end
			if q2054 ~= nil and (q2054:started() or q2054:completed()) then
				me:map(105040314, 0)
			elseif q2053 ~= nil and (q2053:started() or q2053:completed()) then
				me:map(105040312, 0)
			elseif q2052 ~= nil and (q2052:started() or q2052:completed()) then
				me:map(105040310, 0)
			end
		else
			me:dialog(npc, "이상하게 생긴 석상입니다.")
		end
	end
}

-- NPC name (String.wz/Npc.img.xml): 웃고있는고양이

local MAPS = {
	{ map = 801000100, name = "남탕", tail = "#k으" },
	{ map = 801000200, name = "여탕", tail = "#k으" },
}

return {
	on_click = function(me, npc)
		if me:map():wz():id() / 800000000 >= 5 then
			local saved = me:saved_location("WORLDTOUR")
			if saved == nil then
				return
			end
			local sel = me:dialog_list(npc, "안냥? 탈의실 엿보고 싶은거냥?", { "이전에 있던 #m" .. saved .. "# 마을로 돌아갑니다." })
			if sel == nil then
				return
			end
			if not me:dialog_yes_no(npc, "정말 이전에 계시던 #b#m" .. saved .. "##k 마을로 돌아가시고 싶으세요? 돌아가는 요금은 무료이며 저를 통해 언제든지 이곳으로 다시 여행하실 수 있답니다. 지금 돌아가고 싶으세요?") then
				me:dialog(npc, "뭐냥? 쫄았냥? 그정도 담력으론 탈의실 못엿본다냥 저리가라냥!")
				return
			end
			me:clear_saved_location("WORLDTOUR")
			me:map(saved)
			return
		end

		local names = {}
		for i, m in ipairs(MAPS) do
			names[i] = m.name
		end
		local sel = me:dialog_list(npc, "탈의실 엿보고싶냥?", names)
		if sel == nil then
			return
		end
		local selected = MAPS[sel]
		local price = 3000
		if me:class() == Class.Beginner then
			price = 300
		end
		if not me:dialog_yes_no(npc, "정말 #b" .. selected.name .. selected.tail .. "로 가볼테냥? 공짜는 아냥 요금은 #b" .. price .. "#k 메소이며, 구타 당해도 원망하기없기다냥") then
			me:dialog(npc, "여기까지와서 포기할거냥? 실망이다냥.")
			return
		end

		local code = me:exchange({ meso = price }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음. 메소가 부족하신 것 같은데요? 요금이 없으시다면 선택하신 곳으로 보내드릴 수 없답니다.")
			return
		end
		me:save_location("WORLDTOUR")
		me:map(selected.map)
	end
}

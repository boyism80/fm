-- NPC name (String.wz/Npc.img.xml): 콘페이

local MAPS = {
	{ map = 801040004, name = "아래층", tail = "#k으" },
}

return {
	on_click = function(me, npc)
		if me:map():wz():id() / 800000000 >= 5 then
			local saved = me:saved_location("WORLDTOUR")
			if saved == nil then
				return
			end
			local sel = me:dialog_list(npc, "어디로 가겠나?", { "이전에 있던 #m" .. saved .. "# 마을로 돌아갑니다." })
			if sel == nil then
				return
			end
			if not me:dialog_yes_no(npc, "ㅇㅇ #b#m" .. saved .. "##k ㅇㅇ") then
				me:dialog(npc, "...기회 되면 보지.")
				return
			end
			me:map(saved)
			me:clear_saved_location("WORLDTOUR")
			return
		end

		local names = {}
		for i, m in ipairs(MAPS) do
			names[i] = m.name
		end
		local sel = me:dialog_list(npc, "어디로 가겠나, 모험가?", names)
		if sel == nil then
			return
		end
		local selected = MAPS[sel]
		local price = 3000
		if me:class() == Class.Beginner then
			price = 300
		end
		if not me:dialog_yes_no(npc, "정말 #b" .. selected.name .. selected.tail .. "로 갈텐가? 요금은 #b" .. price .. "#k 메소네.") then
			me:dialog(npc, "아직 이곳에 볼일이 남아 있으신가보죠? 여행을 떠나고 싶으시다면 언제든지 제게 다시 찾아와 주세요.")
			return
		end

		local code = me:exchange({ meso = price }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음. 메소가 부족한 것 같은데? 요금이 없으면 할복해.")
			return
		end
		me:save_location("WORLDTOUR")
		me:map(selected.map)
	end
}

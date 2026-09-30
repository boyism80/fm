-- NPC name (String.wz/Npc.img.xml): 스피넬

local PROMPT = "안녕하세요~ 새로운 세계로 여행을 떠나보고 싶으세요? 그렇다면 저희 세계여행 서비스를 이용해보시기 바랍니다! 어디로 여행해보고 싶으세요? 초보자에게는 특별히 90% 할인된 요금으로 제공해드립니다."

local MAPS = {
	{ map = 500000000, name = "태국 플로팅 마켓", tail = "#k으" },
	{ map = 740000000, name = "대만 서문정", tail = "#k으" },
	{ map = 800000000, name = "일본 버섯신사", tail = "#k" },
	{ map = 540010000, name = "싱가포르 중심업무지구", tail = "#k" },
	{ map = 800040000, name = "벚꽃 성 외부지역", tail = "#k" },
	{ map = 701000000, name = "중국 상해와이탄", tail = "#k으" },
}

return {
	on_click = function(me, npc)
		if me:map():wz():id() / 100000000 >= 5 then
			local saved = me:saved_location("WORLDTOUR")
			if saved == nil then
				return
			end
			local sel = me:dialog_list(npc, PROMPT, { "이전에 있던 #m" .. saved .. "# 마을로 돌아갑니다." })
			if sel == nil then
				return
			end
			if not me:dialog_yes_no(npc, "정말 이전에 계시던 #b#m" .. saved .. "##k 마을로 돌아가시고 싶으세요? 돌아가는 요금은 무료이며 저를 통해 언제든지 이곳으로 다시 여행하실 수 있답니다. 지금 돌아가고 싶으세요?") then
				me:dialog(npc, "아직 이곳에 볼일이 남아 있으신가보죠? 여행을 떠나고 싶으시다면 언제든지 제게 다시 찾아와 주세요.")
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
		local sel = me:dialog_list(npc, PROMPT, names)
		if sel == nil then
			return
		end
		local selected = MAPS[sel]
		local price = 3000
		if me:class() == Class.Beginner then
			price = 300
		end
		if not me:dialog_yes_no(npc, "정말 #b" .. selected.name .. selected.tail .. "로 여행을 떠나보시고 싶으세요? 요금은 #b" .. price .. "#k 메소이며, 저를 통해 언제든지 이곳으로 다시 돌아오실 수 있습니다.") then
			me:dialog(npc, "아직 이곳에 볼일이 남아 있으신가보죠? 여행을 떠나고 싶으시다면 언제든지 제게 다시 찾아와 주세요.")
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

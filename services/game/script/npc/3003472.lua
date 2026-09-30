-- NPC name (String.wz/Npc.img.xml): 리안

local DESTINATIONS = {
	{ map = 100000208, name = "제1전장", tail = "#k으" },
	{ map = 130000000, name = "에레브", tail = "#k" },
}

return {
	on_click = function(me, npc)
		local field = me:map()
		local wz = field ~= nil and field:wz() or nil
		if wz == nil then
			return
		end

		if wz:id() >= 500000000 then
			local saved = me:saved_location("WORLDTOUR")
			if saved == nil then
				return
			end
			local sel = me:dialog_list(npc, "전장은 위험해......\r\n#b", {
				" 이전에 있던 #m" .. saved .. "# 마을로 돌아갑니다.",
			})
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

		local options = {}
		for i, dest in ipairs(DESTINATIONS) do
			options[i] = " " .. dest.name
		end
		local sel = me:dialog_list(npc, "전장은 위험해......\r\n#b", options)
		if sel == nil then
			return
		end

		local price = 3000
		if me:class() == Class.Beginner then
			price = 300
		end
		local dest = DESTINATIONS[sel]
		if not me:dialog_yes_no(npc, "정말 #b" .. dest.name .. dest.tail .. "로 여행을 떠나보시고 싶으세요? 요금은 #b" .. price .. "#k 메소이며, 저를 통해 언제든지 이곳으로 다시 돌아오실 수 있습니다.") then
			me:dialog(npc, "아직 이곳에 볼일이 남아 있으신가보죠? 여행을 떠나고 싶으시다면 언제든지 제게 다시 찾아와 주세요.")
			return
		end

		local code = me:exchange({ meso = price }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음. 메소가 부족하신 것 같은데요? 요금이 없으시다면 선택하신 곳으로 보내드릴 수 없답니다.")
			return
		end
		me:save_location("WORLDTOUR")
		me:map(dest.map)
	end
}

-- NPC name (String.wz/Npc.img.xml): 날치

local DESTINATIONS = {
	{ map = 450003000, name = "레헬른", tail = "#k으" },
}

return {
	on_click = function(me, npc)
		local options = {}
		for i, dest in ipairs(DESTINATIONS) do
			options[i] = " " .. dest.name
		end
		local sel = me:dialog_list(npc, "어디로 가시겠어요, 모험가님?\r\n#b", options)
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

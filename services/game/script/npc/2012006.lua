-- NPC name (String.wz/Npc.img.xml): 이스

local platforms = {
	{ map = 200000111, name = "엘리니아", vehicle = "배" },
	{ map = 200000121, name = "루디브리엄", vehicle = "배" },
	{ map = 200000131, name = "리프레", vehicle = "배" },
	{ map = 200000141, name = "무릉", vehicle = "학" },
	{ map = 200000151, name = "아리안트", vehicle = "지니" },
}

return {
	on_click = function(me, npc)
		local labels = {}
		for i, p in ipairs(platforms) do
			labels[i] = p.name .. " 로 가는 " .. p.vehicle .. " 승강장"
		end
		local select = me:dialog_list(npc,
			"오르비스 정거장은 굉장히 넓고 승강장도 비슷비슷해서 길을 잃는 분들이 많답니다. 어느 곳으로 가시려는 건가요?",
			labels)
		if select == nil then
			me:dialog(npc, "목적지를 잘 확인하신 후 저를 통해 승강장으로 이동해 주세요. 시간에 맞춰 서두르시는 게 좋을 거예요!")
			return
		end
		local p = platforms[select]
		if p == nil then
			return
		end
		if not me:dialog_yes_no(npc,
			"혹시 잘못된 승강장으로 가셔도 포탈을 통해 다시 돌아오실 수 있으니 걱정 마세요. #b" ..
				p.name .. " 로 가는 " .. p.vehicle .. " 승강장#k으로 이동하시겠어요?") then
			me:dialog(npc, "목적지를 잘 확인하신 후 저를 통해 승강장으로 이동해 주세요. 시간에 맞춰 서두르시는 게 좋을 거예요!")
			return
		end
		me:map(p.map)
	end,
}

-- NPC name (String.wz/Npc.img.xml): 리린

local EMPRESS_WEAPONS = { 1372084, 1382104, 1402095, 1412065, 1422066, 1432086, 1452111, 1462099, 1472122, 1332130, 1482084, 1492085, 1442116 }
local BASE_WEAPONS = { 1372045, 1382059, 1402047, 1412034, 1422038, 1432049, 1452059, 1462051, 1472071, 1332076, 1482024, 1492025, 1442067 }
local REFUSE = "뭐라구여? 제작하기 싫다고요? 밤길 조심해양!"

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "고렙컨텐츠를 담당하고 있습니다! 골라보세요!#b\r\n", {
			"시간의 돌을 제작한다.",
			"꿈의 돌을 제작한다.",
			"여제 무기를 제작한다",
		})
		if sel == nil then
			me:dialog(npc, REFUSE)
			return
		end

		if sel == 1 then
			if not me:dialog_yes_no(npc, "시간의 조각 10개 가져오면 시간의 돌을 만들어드려요!\r\n물건 없이 말 건거면 밤길 조심하세요!") then
				me:dialog(npc, REFUSE)
				return
			end
			local code = me:exchange({ item = { [4020009] = 10 } }, { item = { [4021010] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:dialog(npc, "기타창 하나 비우세요 칼 맞을수도 있어요.")
				return
			end
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "재료가 이상한데요?")
				return
			end
			me:dialog(npc, "요기 받으세요!")
			return
		end

		if sel == 2 then
			if not me:dialog_yes_no(npc, "꿈의 조각 10개랑 시간의 돌 2개 가져오면 꿈의 돌을 만들어드려요!") then
				me:dialog(npc, REFUSE)
				return
			end
			local code = me:exchange({ item = { [4020013] = 10, [4021010] = 2 } }, { item = { [4021019] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:dialog(npc, "기타창 하나 비우세요 칼 맞을수도 있어요.")
				return
			end
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "재료가 이상한데요?")
				return
			end
			me:dialog(npc, "요기 받으세요!")
			return
		end

		local options = {}
		for i, id in ipairs(EMPRESS_WEAPONS) do
			options[i] = "#i" .. id .. "# #b#z" .. id .. "##k"
		end
		local weapon = me:dialog_list(npc, "만들고 싶은 무기 고르세요!\r\n", options)
		if weapon == nil then
			me:dialog(npc, REFUSE)
			return
		end
		if not me:dialog_yes_no(npc, "#i2049004# 10개 #i4021019# 6개 #i" .. BASE_WEAPONS[weapon] .. "# 1개 필요합니당!") then
			me:dialog(npc, REFUSE)
			return
		end
		local code = me:exchange(
			{ item = { [BASE_WEAPONS[weapon]] = 1, [4021019] = 6, [2049004] = 10 } },
			{ item = { [EMPRESS_WEAPONS[weapon]] = 1 } }
		)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 이상한데요?")
		end
	end
}

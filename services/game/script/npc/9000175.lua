-- NPC name (String.wz/Npc.img.xml): 드리미

local CANCEL = "뭐라구여? 제작하기 싫다고요? 밤길 조심해양!"

local WEAPONS = { 1372222, 1382259, 1402251, 1412177, 1422184, 1432214, 1452252, 1462239, 1472261, 1332274, 1482216, 1492231, 1442268 }
local MATERIALS = { 1372084, 1382104, 1402095, 1412065, 1422066, 1432086, 1452111, 1462099, 1472122, 1332130, 1482084, 1492085, 1442116 }

return {
	on_click = function(me, npc)
		local menu = me:dialog_list(npc, "반짝이는 물건을 좋아한다옹~", {
			"앱솔 랩스 코인을 제작한다.",
			"앱솔 랩스 무기를 제작한다.",
		})
		if menu == nil then
			me:dialog(npc, CANCEL)
			return
		end

		if menu == 1 then
			if not me:dialog_yes_no(npc, "#i3980091##k네잎 클로버 코인 5개랑 #i3980090##k정화의 조각 5000개를 가져오면 앱솔랩스 코인을 주겠다냥.") then
				me:dialog(npc, CANCEL)
				return
			end
			local code = me:exchange({ item = { [3980090] = 5000, [3980091] = 5 } }, { item = { [3980092] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:dialog(npc, "기타창 하나 비우세요 칼 맞을수도 있어요.")
				return
			elseif code ~= ExchangeResult.OK then
				me:dialog(npc, "재료가 이상한데요?")
				return
			end
			me:dialog(npc, "요기 받으세요!")
			return
		end

		local options = {}
		for i, item in ipairs(WEAPONS) do
			options[i] = "#i" .. item .. "# #b#z" .. item .. "##k"
		end
		local sel = me:dialog_list(npc, "만들고 싶은 무기 고르세요!", options)
		if sel == nil then
			me:dialog(npc, CANCEL)
			return
		end
		if not me:dialog_yes_no(npc, "#i3980092# 1개 #i3980090# 500개 #i" .. MATERIALS[sel] .. "# 1개 필요합니당!") then
			me:dialog(npc, CANCEL)
			return
		end

		local code = me:exchange({ item = { [MATERIALS[sel]] = 1, [3980092] = 1, [3980090] = 500 } }, { item = { [WEAPONS[sel]] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 이상한데요?")
		end
	end
}

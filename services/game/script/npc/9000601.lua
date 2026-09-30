-- NPC name (String.wz/Npc.img.xml): 오르카

local BUNGEOPPANG = 3980055
local TIERS = {
	{ upto = 10, items = { 5060002, 1142363, 2040811, 5120010, 1012073, 1012072, 1012071, 1012070, 2048010, 2048012, 2048013, 5120010, 5120010, 5120010, 5120010, 5120010 }, announce = true },
	{ upto = 30, items = { 5060002, 2022038, 4005004, 4021007, 4011002, 4011006, 2022038, 2022038, 2022038, 2022038, 2022038, 2022038 } },
	{ upto = 60, items = { 5060002, 2022038, 5510000 } },
	{ upto = 100, items = { 5060002, 2022038 } },
}

return {
	on_click = function(me, npc)
		local message = "#i" .. BUNGEOPPANG .. "# #b#z" .. BUNGEOPPANG .. "##k으로 진귀한 아이템을 교환해드려요! \r\n #i" .. BUNGEOPPANG .. "# #b#z" .. BUNGEOPPANG .. "##k는 퀘스트로 획득합니다.\r\n"
		local sel = me:dialog_list(npc, message, {
			"#b#z" .. BUNGEOPPANG .. "#을 사용한다.",
		})
		if sel == nil then
			return
		end

		local roll = math.random(1, 100)
		local tier = nil
		for _, t in ipairs(TIERS) do
			if roll <= t.upto then
				tier = t
				break
			end
		end
		local item_id = tier.items[math.random(1, #tier.items)]
		if me:exchange({ item = { [BUNGEOPPANG] = 1 } }, { item = { [item_id] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 빈 공간이 없거나 #b붕어빵#k이 부족합니다.퀘스트를 클리어하시거나 베르가모트 에게서 구하세요.")
			return
		end
		me:dialog(npc, "#i" .. item_id .. ":# #b#z" .. item_id .. "##k이(가) 나왔습니다.")
		if tier.announce then
			me:world_message(2, "[이벤트] : " .. me:name() .. " 님이 10% 확률로 진귀한 아이템을 획득하였습니다.")
		end
	end
}

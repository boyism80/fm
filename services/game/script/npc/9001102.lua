-- NPC name (String.wz/Npc.img.xml): 아기 월묘

local MAPLE_LEAF = 3980000
local WHITE_SCROLL = 2049004
local NECKLACE = 1122014
local DRAGON_EGG = 4001094
local STONE = 2041200
local FEE = 10000

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "안녕하세요! 저는 #b목걸이 컨텐츠#k를 담당하는 아기월묘 라고 합니다! 원하는걸 골라보세요!\r\n", {
			" #r드래곤의 알#k을 가져왔어요!.",
			" #b영주목걸이#k를 각성 시키고 싶어요!",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			if me:dialog_yes_no(npc, "#i" .. STONE .. "# 을 만들려면\r\n#i" .. DRAGON_EGG .. "# 1개와 메소 만원이필요합니다.") == false then
				return
			end
			if item_count(me, DRAGON_EGG) < 1 or me:meso() < FEE then
				me:dialog(npc, "재료가부족해요ㅠㅠㅠ")
				return
			end
			if me:exchange({ item = { [DRAGON_EGG] = 1 }, meso = FEE }, { item = { [STONE] = 1 } }) ~= ExchangeResult.OK then
				me:dialog(npc, "아이템창 1칸 부족해엿")
				return
			end
			me:dialog(npc, "완성되었어요!")
			return
		end

		if me:dialog_yes_no(npc, "#i" .. NECKLACE .. "# (공+5 마+10 올텟11) 만들려면\r\n#i" .. NECKLACE .. "# 1개, #i" .. MAPLE_LEAF .. "# 1200개와 백줌 100퍼 1개가 필요합니다.") == false then
			return
		end
		if item_count(me, NECKLACE) < 1 or item_count(me, WHITE_SCROLL) < 1 or item_count(me, MAPLE_LEAF) < 1200 then
			me:dialog(npc, "재료가부족해요ㅋ큐ㅠㅠㅠ")
			return
		end
		local cost = { item = { [NECKLACE] = 1, [WHITE_SCROLL] = 1, [MAPLE_LEAF] = 1200 } }
		local reward = {
			item = { [NECKLACE] = 1 },
			bonus = { str = 10, dex = 10, int = 10, luk = 10, watk = 5, matk = 10 },
		}
		if me:exchange(cost, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "아이템창 1칸 부족해엿")
			return
		end
		me:dialog(npc, "완성되었어요!")
	end
}

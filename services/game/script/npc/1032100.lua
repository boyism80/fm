-- NPC name (String.wz/Npc.img.xml): 요정 아르웬

local RECIPES = {
	{
		name = "달의 돌",
		item = 4011007,
		cost = 10000,
		mats = { 4011000, 4011001, 4011002, 4011003, 4011004, 4011005, 4011006 },
	},
	{
		name = "별의 돌",
		item = 4021009,
		cost = 15000,
		mats = { 4021000, 4021001, 4021002, 4021003, 4021004, 4021005, 4021006, 4021007, 4021008 },
	},
	{
		name = "검은 깃털",
		item = 4031042,
		cost = 30000,
		mats = { 4011007, 4001006, 4021008 },
	},
}

return {
	on_click = function(me, npc)
		if me:level() < 40 then
			me:dialog(npc, "저만이 만들 수 있는 희귀한 물건은 있지만.. 아직 약하시고.. 누군지도 모르는 당신에게는 만들어 드릴 수 없어요.")
			return
		end
		if not me:dialog(npc, "음.. 제가 바로 소문의 최고의 연금술사입니다. 음.. 오랜 시간동안 인간과 접촉한 요정은 없었지만.. 당신같이 강한 사람은 괜찮아 보이는군요.. 특별히 당신을 위해서는 아이템을 제작해 드리도록 하지요.", false, true) then
			return
		end

		local sel = me:dialog_list(npc, "무얼 만들어 보고 싶으신가요?\r\n#b", {
			"달의 돌",
			"별의 돌",
			"검은 깃털",
		})
		if sel == nil then
			return
		end

		local recipe = RECIPES[sel]
		local text = "만들고 싶은 아이템이 #b#t" .. recipe.item .. "##k 인가요? 재료는 다음과 같아요.\r\n"
		for _, mat in ipairs(recipe.mats) do
			text = text .. "\r\n#b#i" .. mat .. "# #t" .. mat .. "# 1 개"
		end
		text = text .. "\r\n#b#i4031138# " .. recipe.cost .. " 메소"
		if not me:dialog_yes_no(npc, text) then
			me:dialog(npc, recipe.name .. "을 제작하는건 쉬운일이 아니랍니다.. 재료를 제게 가져와 주세요.", false, true)
			return
		end

		local fail = "메소는 충분히 갖고 계신지, 또는 재료가 부족한건 아닌지, 인벤토리 공간이 충분한지 다시 한번 확인해 주세요."
		if me:meso() <= recipe.cost then
			me:dialog(npc, fail, false, true)
			return
		end
		local cost_items = {}
		for _, mat in ipairs(recipe.mats) do
			cost_items[mat] = 1
		end
		local code = me:exchange({ item = cost_items, meso = recipe.cost }, { item = { [recipe.item] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, fail, false, true)
			return
		end
		me:dialog(npc, "다 되었답니다. 더 필요한 물건이 있으시면 언제라도 다시 찾아와주세요..", false, true)
	end
}

-- NPC name (String.wz/Npc.img.xml): 체프

local SCROLL = 4031048
local rewards = {
	{ 2000004, 100 },
	{ 4011006, 10 },
	{ 4011000, 10 },
	{ 4011005, 10 },
	{ 4021005, 10 },
	{ 4021001, 10 },
	{ 4021007, 10 },
}

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		if item_count(me, SCROLL) < 1 then
			return
		end
		if not me:dialog_yes_no(npc, "어? 그건 #b#t4031048##k 맞지? 음.. 그것을 나에게 주지 않겠어? 보상은 섭섭치 않게 해주지.") then
			return
		end
		local reward = rewards[math.random(1, #rewards)]
		local code = me:exchange({ item = { [SCROLL] = 1 } }, { item = { [reward[1]] = reward[2] } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 충분한지, 혹은 비밀의 주문서는 제대로 갖고 있는거야?")
			return
		end
		me:dialog(npc, "자, #b#t" .. reward[1] .. "# " .. reward[2] .. "#k개를 줄게. 고마워~")
	end
}

-- NPC name (String.wz/Npc.img.xml): 몽

local MAPLE_LEAF = 4001126
local COST = 100

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		if me:dialog_yes_no(npc, "단풍잎 100개를 이용하여 #r#bMP -100#k을 내려드립니다.\r\n") == false then
			return
		end

		if me:base_hp() < 2000 then
			me:dialog(npc, "피가 너무 적어요.")
			return
		end
		if me:base_mp() < 2001 then
			me:dialog(npc, "마나가 너무 많아요.")
			return
		end
		if item_count(me, MAPLE_LEAF) < COST then
			me:dialog(npc, "단풍잎이 부족 합니다.")
			return
		end
		if me:exchange({ item = { [MAPLE_LEAF] = COST } }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "단풍잎이 부족 합니다.")
			return
		end
		me:max_mp(me:base_mp() - 100)
	end
}

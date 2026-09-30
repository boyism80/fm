-- NPC name (String.wz/Npc.img.xml): 핑크몽

local JOURNAL = 3980000
local COST = 60

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		if me:dialog_yes_no(npc, "퀘스트일지 60개를 이용하여 #rHP 10#k, #bMP 120#k을 올려드립니다.\r\n") == false then
			return
		end

		if me:base_hp() > 29500 then
			me:dialog(npc, "피가 너무 많아요.")
			return
		end
		if me:base_mp() > 29900 then
			me:dialog(npc, "마나가 너무 많아요.")
			return
		end
		if item_count(me, JOURNAL) < COST then
			me:dialog(npc, "퀘스트일지가 부족합니다.")
			return
		end
		if me:exchange({ item = { [JOURNAL] = COST } }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "퀘스트일지가 부족합니다.")
			return
		end
		me:max_hp(me:base_hp() + 10)
		me:max_mp(me:base_mp() + 120)
	end
}

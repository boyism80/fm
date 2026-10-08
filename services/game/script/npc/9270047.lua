-- NPC name (String.wz/Npc.img.xml): 알돌

local COIN = 3980018

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "환상의 나라 최강의 무투가 타르가 에게 도전하시겠어요??", { "타르가와 붙어본다.." })
		if sel ~= 1 then
			return
		end
		if item_count(me, COIN) < 1 then
			me:dialog(npc, "타르가 코인 1개가 필요합니다.")
			return
		end
		if me:records():get("summon.targa") >= 10 then
			me:dialog(npc, "오늘 이미 소환하셨네요? 다음에 시도해주세요.~")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "맵 안에 아직 몬스터 있습니다.~")
			return
		end

		local code = me:exchange({ item = { [COIN] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "타르가 코인 1개가 필요합니다.")
			return
		end
		me:records():add("summon.targa", 1, { daily = true })
		map:spawn_mob(9420541, -397, 640)
	end
}

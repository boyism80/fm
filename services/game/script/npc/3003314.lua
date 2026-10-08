-- NPC name (String.wz/Npc.img.xml): 리더 알

local SPAWN_X = -963
local SPAWN_Y = 165
local BUTTERFLY = 3980009
local DAILY_LIMIT = 10

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "정말 위험할 것 같은데요?#b", {
			"분노한 거대 정령과 붙어보겠다..",
		})
		if sel == nil then
			return
		end

		if item_count(me, BUTTERFLY) < 1000 then
			me:dialog(npc, "루시드의 나비 1000개가 필요합니다.")
			return
		end
		local count = me:records():get("summon.corrupted_spirit_of_harmony")
		if count >= DAILY_LIMIT then
			me:dialog(npc, "오늘 이미 소환하셨네요? 다음에 시도해주세요.~")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "맵 안에 아직 몬스터 있습니다.~")
			return
		end
		if me:exchange({ item = { [BUTTERFLY] = 1000 } }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "루시드의 나비 1000개가 필요합니다.")
			return
		end
		me:records():add("summon.corrupted_spirit_of_harmony", 1, { daily = true })
		map:spawn_mob(8644011, SPAWN_X, SPAWN_Y)
	end
}

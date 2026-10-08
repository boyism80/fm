-- NPC name (String.wz/Npc.img.xml): 아미

local SPAWN_X = -10
local SPAWN_Y = -215
local HELMET = 1003112
local ZAKUM_PIECE = 4033329
local CHAOS_PIECE = 4033330
local CHAOS_ZAKUM_ARMS = { 8800103, 8800104, 8800105, 8800106, 8800107, 8800108, 8800109, 8800110 }

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "카오스 자쿰을 소환하실건가요?#b", {
			"카오스 자쿰을 소환한다.",
			"카오스 자쿰의 투구 각성.",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			if item_count(me, ZAKUM_PIECE) < 10 then
				me:dialog(npc, "자쿰의 결정 10개가 필요합니다.")
				return
			end
			if me:records():get("summon.chaos_zakum") >= 1 then
				me:dialog(npc, "오늘 이미 소환하셨네요? 다음에 시도해주세요.~")
				return
			end
			local map = me:map()
			if next(map:mobs()) ~= nil then
				me:dialog(npc, "맵 안에 아직 몬스터 있습니다.~")
				return
			end
			if me:exchange({ item = { [ZAKUM_PIECE] = 10 } }, nil) ~= ExchangeResult.OK then
				me:dialog(npc, "자쿰의 결정 10개가 필요합니다.")
				return
			end
			me:records():add("summon.chaos_zakum", 1, { daily = true })
			map:spawn_mob(8800100, SPAWN_X, SPAWN_Y)
			for _, mob_id in ipairs(CHAOS_ZAKUM_ARMS) do
				map:spawn_mob(mob_id, SPAWN_X, SPAWN_Y)
			end
			return
		end

		if me:records():get("chaos_zakum.helmet_upgrade") >= 1 then
			me:dialog(npc, "이미 모자를 한번 강화하셨어요.")
			return
		end
		if item_count(me, HELMET) < 1 or item_count(me, CHAOS_PIECE) < 5 then
			me:dialog(npc, "아이템이 부족합니다. 재료는 카쿰투구와 카쿰결정 5개입니다. 주문서작을 하기전에 강화 해주세요.")
			return
		end
		local cost = { item = { [HELMET] = 1, [CHAOS_PIECE] = 5 } }
		local reward = {
			item = { [HELMET] = 1 },
			bonus = { str = 5, dex = 5, int = 5, luk = 5, watk = 5, matk = 5 },
		}
		if me:exchange(cost, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "장비창 한칸을 비우시오.")
			return
		end
		me:records():set("chaos_zakum.helmet_upgrade", 1)
	end
}

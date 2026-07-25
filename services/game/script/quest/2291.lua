-- VIP존 입장하기 (Quest.wz/QuestData/2291.img): VIP존 입장하기

local quest_id = 2291
local vip_base = 103040440

local player_count, reset_vip_map, on_end

local function map_player_count(map_id)
	local ok, count = run_on_map(map_id, "script/quest/2291.lua", "player_count")
	if not ok then
		return -1
	end
	if count == nil then
		return 0
	end
	return count
end

local function reset_vip_maps(map_id)
	run_on_map(map_id, "script/quest/2291.lua", "reset_vip_map")
end

return {
	player_count = function(map)
		local n = 0
		for _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	reset_vip_map = function(map)
		map:reset()
		return true
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_yes_no(npc, "소지하신 VIP티켓은 모두 모아 온건가요? 어서 제게 주세요. #bVIP존#k으로 입장시켜 드릴게요. 참고로 VIP존은 #b30분에 한 번씩만 이용이 가능#k합니다.") then
			me:dialog(npc, "지금은 별로 가고 싶지 않으신가 보군요. 언제든 저를 통해서 입장할 수 있으니 말을 걸어주세요.", false, false)
			return
		end

		for i = 0, 9 do
			local empty = true
			local ids = { vip_base + i, vip_base + 10 + i, vip_base + 20 + i }
			for _, map_id in ipairs(ids) do
				if map_player_count(map_id) ~= 0 then
					empty = false
					break
				end
			end
			if empty then
				for _, map_id in ipairs(ids) do
					reset_vip_maps(map_id)
				end
				local code = me:exchange({ item = { [4032521] = 10 } }, {})
				if code ~= ExchangeResult.OK then
					return
				end
				q:force_complete(npc)
				me:show_effect(EffectType.QuestCompletion)
				me:play_portal_sound()
				me:map(vip_base + i, 1)
				return
			end
		end

		me:dialog(npc, "모든 VIP존에 누군가 있군요. 다른 채널을 이용해주세요.", false, false)
	end
}

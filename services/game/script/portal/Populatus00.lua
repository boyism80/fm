-- Portal (old/scripts/portal/Populatus00.js): 시계탑 깊은 곳 → 파풀라투스

local pq = require("script/lib/party_quest")

local MEDAL_ID = 4031172
local QUEST_MAP = 922020300
local BOSS_MAP = 220080001
local INSTANCE_ID = "Battle"
local MAX_PLAYERS = 12

local function quest_started(me, id)
	local q = me:quest(id)
	return q ~= nil and q:started()
end

local function quest_record(me, id)
	local q = me:quest(id)
	if q == nil then
		return nil
	end
	return q:record()
end

local function player_count(map)
	if map == nil then
		return 0
	end
	local n = 0
	for _, ch in pairs(map:characters()) do
		if ch ~= nil then
			n = n + 1
		end
	end
	return n
end

return {
	on_enter = function(me)
		if not pq.has_item(me, MEDAL_ID) then
			me:notice("<루디브리엄의 메달>을 갖고 있어야 이 문을 지날 수 있습니다.")
			return
		end
		local record6364 = quest_record(me, 6364)
		if quest_started(me, 6361)
			or quest_started(me, 6362)
			or (quest_started(me, 6363) and (record6364 == nil or record6364 == "")) then
			me:play_portal_sound()
			me:map(QUEST_MAP)
			return
		end
		local group = state_machine("papulatus")
		if group == nil then
			return
		end
		if group:get_property("battle") == "1" then
			me:notice("이미 파풀라투스와의 전투가 시작되어 입장할 수 없습니다.")
			return
		end
		local boss_map = group:map(BOSS_MAP)
		if boss_map == nil then
			boss_map = id2map(BOSS_MAP)
		end
		if player_count(boss_map) > MAX_PLAYERS then
			me:notice("이 방은 이미 파풀라투스와의 전투를 위한 최대 인원수 만큼 가득 찼습니다.")
			return
		end
		me:play_portal_sound()
		local sm = group:get(INSTANCE_ID)
		if sm == nil then
			local err
			sm, err = group:create(INSTANCE_ID, me)
			if sm == nil then
				if err ~= nil then
					log("papulatus create:", err)
				end
				return
			end
		end
		sm:enter_player(me)
	end
}

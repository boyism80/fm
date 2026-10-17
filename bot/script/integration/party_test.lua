local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local PARK = 100000200
local HUNTING_GROUND = 104040000
local SNAIL = 100100
local SNAIL_QUEST = 1019

local STATUS = { different_map = 27, different_channel = 29 }

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot)
	if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function change_leader_refused(ctx, leader, target, code, what)
	local p = leader:request(resp.party_status_message, req.party_operation { operation = PARTY.ChangeLeader, target_character_id = target:id() }, function(p)
		return p.code == code
	end, 5000)
	return check(ctx, p ~= false, string.format("%s 파티장 위임이 거절되지 않음 (기대 코드 %d)", what, code))
end

-- The leader can hand over only to a member on the same channel and map, and a leader who leaves disbands the party.
local function leader_flow(ctx)
	local a, b, c = ctx:bot(0), ctx:bot(1), ctx:bot(2)
	for _, bot in ipairs({ a, b, c }) do
		if setup(ctx, bot) == false then
			return false
		end
	end
	if pq.form_party(ctx) == false then
		return false
	end

	if b:map_move(PARK) == false then
		return ctx:fail("다른 맵 이동 실패")
	end
	ctx:sleep(1000)
	if change_leader_refused(ctx, a, b, STATUS.different_map, "다른 맵") == false then
		return false
	end

	if b:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	ctx:sleep(1000)
	if change_leader_refused(ctx, a, b, STATUS.different_channel, "다른 채널") == false then
		return false
	end
	if b:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	if b:map_move(HENESYS) == false then
		return ctx:fail("같은 맵 복귀 실패")
	end
	ctx:sleep(1000)

	local changed = a:request_on(c, resp.party_update_leader_change, req.party_operation { operation = PARTY.ChangeLeader, target_character_id = b:id() }, nil, 5000)
	if check(ctx, changed ~= false and changed.new_leader_character_id == b:id(), "같은 맵 파티장 위임이 다른 파티원에게 보이지 않음") == false then
		return false
	end

	local left = c:request_on(a, resp.party_update_leave, req.party_operation { operation = PARTY.Leave }, function(p)
		return p.target_character_id == c:id()
	end, 5000)
	if check(ctx, left ~= false and left.leader_character_id == b:id(), "파티원 탈퇴가 남은 파티원에게 보이지 않음") == false then
		return false
	end

	local disband = b:request_on(a, resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	if check(ctx, disband ~= false, "파티장이 나갔는데 파티가 해체되지 않음") == false then
		return false
	end
	ctx:sleep(500)

	if a:request(resp.party_created, req.party_operation { operation = PARTY.Create }, nil, 5000) == false then
		return ctx:fail("해체된 파티에 남아 새 파티를 만들 수 없음")
	end
	if a:request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000) == false then
		return ctx:fail("새 파티 해체 실패")
	end
	return true
end

local function spawn_snail(ctx, bot)
	local before = {}
	for _, mob in ipairs(bot:mobs(SNAIL)) do
		before[mob.oid] = true
	end
	if bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL }) == false then
		return ctx:fail("달팽이 생성 응답 없음")
	end
	for _, mob in ipairs(bot:mobs(SNAIL)) do
		if before[mob.oid] == nil then
			return mob.oid
		end
	end
	return ctx:fail("생성한 달팽이를 찾지 못함")
end

-- A party member on the map counts a hunting quest kill even when another member lands the last hit.
local function shared_kill_count(ctx)
	local a, b = ctx:bot(0), ctx:bot(1)
	for _, bot in ipairs({ a, b, ctx:bot(2) }) do
		if pq.command(bot, "/봇초기화 30 0 0 - " .. SNAIL_QUEST .. ":1", "봇초기화 완료") == false then
			return ctx:fail(bot:name() .. " 봇 초기화 실패")
		end
	end
	if pq.form_party(ctx) == false then
		return false
	end
	for _, bot in ipairs({ a, b }) do
		if bot:force_instance_move(HUNTING_GROUND, "파티 처치 수") == false then
			return ctx:fail(bot:name() .. " 사냥터 이동 실패")
		end
		if pq.command(bot, "/즉사 on", "즉사 상태: enabled") == false then
			return ctx:fail(bot:name() .. " 즉사 설정 실패")
		end
	end

	local counted = function(p)
		return p.quest_status.quest_id == SNAIL_QUEST and p.quest_status.status == 1
	end
	for _, pair in ipairs({ { a, b }, { b, a } }) do
		local killer, member = pair[1], pair[2]
		local oid = spawn_snail(ctx, killer)
		if oid == false then
			return false
		end
		if killer:attack(oid, 1) == false then
			return ctx:fail(killer:name() .. " 공격 실패")
		end
		if member:request(resp.update_quest, nil, counted, 5000) == false then
			return ctx:fail(killer:name() .. "가 잡은 달팽이가 " .. member:name() .. "의 사냥 퀘스트 처치 수에 오르지 않음")
		end
	end

	for _, bot in ipairs({ a, b, ctx:bot(2) }) do
		pq.command(bot, "/즉사 off", "즉사 상태: disabled")
		bot:request(resp.party_update_leave, req.party_operation { operation = PARTY.Leave }, nil, 1000)
	end
	return true
end

test_suite {
	name = "Party: 파티장 위임 조건·탈퇴·파티장 탈퇴 해체·처치 수 공유",
	bot_count = 3,

	scenarios = {
		leader_flow,
		shared_kill_count,
	},
}

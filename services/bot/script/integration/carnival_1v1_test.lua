local pq = require("script/integration/lib/party_quest")

local HUB = 980000000
local SPIEGELMANN = 2042000
local REWARD_NPC = 2042002
local GUARDIAN_REACTOR = 9980000
local TEAM_BLUE = 1
local TAB_MOB = 0
local TAB_GUARDIAN = 2
local GUARDIAN_CP = 17
local MOB_CP = 7
local DIALOG_DEFAULT = 0
local DIALOG_YES_NO = 1

local field_number = nil

local function waiting_map()
	return HUB + field_number * 100
end

local function wait_map(ctx, bot, map_id, timeout)
	if bot:map() == map_id then
		return true
	end
	local warp = bot:request(resp.warp, nil, function(p)
		return p.character.map == map_id
	end, timeout)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 없음: " .. map_id .. " (현재 " .. bot:map() .. ")")
	end
	return true
end

local function open_field_list(ctx, bot)
	local oid = pq.npc(ctx, bot, SPIEGELMANN)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail(bot:name() .. " 슈피겔만 대화 응답 없음")
	end
	if dlg.selections == nil then
		bot:dialog(false)
		return ctx:fail(bot:name() .. " 필드 목록 대신 다른 대화: " .. dlg.text)
	end
	return dlg
end

local function earn_cp(ctx, bot, enemy, need)
	while bot:cp() < need do
		local mob = nil
		for _, m in ipairs(bot:mobs()) do
			if m.team == enemy then
				mob = m
				break
			end
		end
		if mob == nil then
			local spawn = bot:request(resp.spawn_mob, nil, function(p)
				return p.mob.carnival_team == enemy
			end, 30000)
			if spawn == false then
				return ctx:fail("상대팀 몹 재생성 대기 시간 초과")
			end
		elseif bot:kill(mob.oid) == false then
			return ctx:fail("상대팀 몹 처치 실패: " .. mob.id)
		else
			bot:request(resp.carnival_obtained_cp, nil, nil, 1000)
		end
	end
	return true
end

test_suite {
	name = "Monster Carnival 1:1",
	bot_count = 2,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 30", "레벨 설정") == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
			end
			if bot:map_move(HUB) == false then
				return false
			end
			if bot:request(resp.party_created, req.party_operation { operation = PARTY.Create }) == false then
				return ctx:fail(bot:name() .. " 파티 생성 실패")
			end
		end
		return true
	end,

	on_finished = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			ctx:bot(i):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
		end
	end,

	scenarios = {
		function(ctx)
			local red = ctx:bot(0)
			local dlg = open_field_list(ctx, red)
			if dlg == false then
				return false
			end
			local index = nil
			for i, label in ipairs(dlg.selections) do
				if label:find("대기", 1, true) == nil then
					index = i
					field_number = tonumber(label:match("%d+"))
					break
				end
			end
			if index == nil then
				return ctx:fail("빈 카니발 필드가 없음")
			end

			if red:dialog(true, index - 1) == nil then
				return ctx:fail("필드 개설 확인 대화 없음")
			end
			local warp = red:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
				return p.character.map == waiting_map()
			end)
			if warp == false then
				return ctx:fail("레드팀 대기실 이동 실패")
			end
			return true
		end,
		function(ctx)
			local red = ctx:bot(0)
			local blue = ctx:bot(1)
			local dlg = open_field_list(ctx, blue)
			if dlg == false then
				return false
			end
			local index = nil
			for i, label in ipairs(dlg.selections) do
				if label:find("카니발 필드" .. field_number .. " (", 1, true) ~= nil then
					index = i
				end
			end
			if index == nil then
				return ctx:fail("대기 중인 필드가 목록에 없음: " .. field_number)
			end

			if blue:dialog(true, index - 1) == nil then
				return ctx:fail("도전 확인 대화 없음")
			end
			local ask = blue:request_on(red, resp.dialog_yes_no, req.dialog { dialog_type = DIALOG_YES_NO, next = true })
			if ask == false then
				return ctx:fail("레드팀 파티장에게 도전 수락 대화가 오지 않음")
			end
			blue:send(req.dialog { dialog_type = DIALOG_DEFAULT, next = false })

			local enter = red:request_on(blue, resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
				return p.character.map == waiting_map()
			end)
			if enter == false then
				return ctx:fail("블루팀 대기실 이동 실패")
			end

			local field = waiting_map() + 1
			if wait_map(ctx, red, field, 20000) == false then
				return false
			end
			return wait_map(ctx, blue, field, 20000)
		end,
		function(ctx)
			local red = ctx:bot(0)
			local blue = ctx:bot(1)
			if earn_cp(ctx, red, TEAM_BLUE, GUARDIAN_CP + MOB_CP) == false then
				return false
			end

			local summon = red:request_on(blue, resp.spawn_reactor, req.carnival { tab = TAB_GUARDIAN, num = 0 }, function(p)
				return p.reactor.reactor_id == GUARDIAN_REACTOR
			end)
			if summon == false then
				return ctx:fail("수호물 소환 실패 (CP " .. red:cp() .. ")")
			end
			local guardian = summon.reactor
			if pq.move(blue, guardian.position.x, guardian.position.y) == false then
				return ctx:fail("수호물 앞으로 이동 실패")
			end
			local destroyed = false
			for _ = 1, 10 do
				local p, name = blue:hit_reactor(guardian.oid)
				if p == false then
					return ctx:fail("수호물 타격 응답 없음")
				end
				if name == resp.destroy_reactor then
					destroyed = true
					break
				end
			end
			if destroyed == false then
				return ctx:fail("10번 때려도 수호물이 부서지지 않음")
			end

			local spawn = red:request_on(blue, resp.spawn_mob, req.carnival { tab = TAB_MOB, num = 0 })
			if spawn == false then
				return ctx:fail("몹 소환 실패 (CP " .. red:cp() .. ")")
			end
			if blue:kill(spawn.mob.oid) == false then
				return ctx:fail("소환된 몹 처치 실패")
			end
			return true
		end,
		function(ctx)
			local red = ctx:bot(0)
			local blue = ctx:bot(1)
			if pq.command(red, "/타이머 1", "남은 시간") == false then
				return ctx:fail("전투 시간 단축 실패")
			end
			local win = waiting_map() + 3
			local lose = waiting_map() + 4
			if wait_map(ctx, red, win, 30000) == false then
				return false
			end
			return wait_map(ctx, blue, lose, 30000)
		end,
		function(ctx)
			local expected = { "축하하네", "아쉽게도" }
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				local oid = pq.npc(ctx, bot, REWARD_NPC)
				if oid == false then
					return false
				end
				local dlg = bot:npc_click(oid)
				if dlg == false then
					return ctx:fail(bot:name() .. " 보상 대화 없음")
				end
				if dlg.text:find(expected[i + 1], 1, true) == nil then
					return ctx:fail(bot:name() .. " 예상과 다른 결과: " .. dlg.text)
				end
				local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == HUB
				end)
				if warp == false then
					return ctx:fail(bot:name() .. " 허브로 돌아가지 못함")
				end
			end
			return true
		end,
	},
}

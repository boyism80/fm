local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")

local HENESYS = 100000000
local ELLINIA = 101000000
local TELEPORT = 0
local COST = 300
local SETTLE_WAIT_MS = 40000

local FAMILY = {
	family_result = true,
	family_invite = true,
	family_invite_result = true,
	family_joined = true,
	family_info = true,
	family_chart = true,
}

local inbox = {}

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function listen(ctx)
	ctx:hook("family_recovery", function(_, bot, p, name)
		if FAMILY[name] == nil then
			return
		end
		local box = inbox[bot:id()]
		if box == nil then
			box = {}
			inbox[bot:id()] = box
		end
		box[#box + 1] = { name = name, packet = p }
	end)
end

local function send(bot, request)
	inbox[bot:id()] = {}
	return bot:send(request)
end

local function expect(ctx, bot, name, match, what)
	for _ = 1, 50 do
		local box = inbox[bot:id()] or {}
		for i, item in ipairs(box) do
			if item.name == name and (match == nil or match(item.packet)) then
				table.remove(box, i)
				return item.packet
			end
		end
		ctx:sleep(100)
	end
	return ctx:fail(string.format("%s %s 응답 없음", bot:name(), what))
end

local function setup(ctx, bot, level, map)
	if pq.command(bot, string.format("/봇초기화 %d 0 0 -", level), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= map and bot:map_move(map) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

-- Internal's own record, read through the pedigree, so a stale family window cannot hide a held cost.
local function reputation(ctx, bot)
	send(bot, req.family_pedigree { family = { name = bot:name() } })
	local chart = expect(ctx, bot, "family_chart", nil, "가계도")
	if chart == false then
		return nil
	end
	for _, member in ipairs(chart.state.members) do
		if member.character_id == bot:id() then
			return member.reputation
		end
	end
	ctx:fail(bot:name() .. " 가계도에 자신이 없음")
	return nil
end

local function holds(ctx, bot, expected, what)
	local have = reputation(ctx, bot)
	if have == nil then
		return false
	end
	return check(ctx, have == expected, string.format("%s %s: 명성 %d (기대 %d)", bot:name(), what, have, expected))
end

local function wait_holds(ctx, bot, expected, what)
	for _ = 1, SETTLE_WAIT_MS / 2000 do
		if reputation(ctx, bot) == expected then
			return true
		end
		ctx:sleep(2000)
	end
	return holds(ctx, bot, expected, what)
end

local function join(ctx, senior, junior)
	send(senior, req.family_remove_senior {})
	send(junior, req.family_remove_senior {})
	ctx:sleep(1000)
	if junior:map_move(ELLINIA) == false then
		return ctx:fail("주니어 맵 이동 실패")
	end
	if senior:map_move(ELLINIA) == false then
		return ctx:fail("시니어 맵 이동 실패")
	end

	send(senior, req.family_invite { family = { name = junior:name() } })
	if expect(ctx, junior, "family_invite", nil, "패밀리 초대 수신") == false then
		return false
	end
	send(junior, req.family_accept { family = { target_id = senior:id(), name = senior:name(), accepted = true } })
	if expect(ctx, junior, "family_joined", nil, "패밀리 가입 알림") == false then
		return false
	end
	if junior:map_move(HENESYS) == false then
		return ctx:fail("주니어 헤네시스 이동 실패")
	end
	return true
end

-- The junior's progress is credited to the senior on save.
local function earn(ctx, senior, junior)
	if pq.command(junior, "/경험치얻기 200000", "경험치 200000 획득") == false then
		return ctx:fail("경험치 획득 실패")
	end
	if pq.command(junior, "/서버저장", "서버 저장이 완료되었습니다") == false then
		return ctx:fail("서버 저장 실패")
	end
	for _ = 1, 20 do
		local have = reputation(ctx, senior)
		if have ~= nil and have >= COST then
			return have
		end
		ctx:sleep(1000)
	end
	ctx:fail("시니어 명성이 쌓이지 않음")
	return nil
end

local function prepare(ctx)
	local senior, junior = ctx:bot(0), ctx:bot(1)
	listen(ctx)
	if setup(ctx, senior, 50, ELLINIA) == false or setup(ctx, junior, 40, ELLINIA) == false then
		return nil
	end
	if join(ctx, senior, junior) == false then
		return nil
	end
	return earn(ctx, senior, junior)
end

local function teleport(senior, junior)
	return req.family_use { family = { feature = TELEPORT, name = junior:name() } }
end

-- The fault is server wide, so it is turned off right after the one request it targets.
local function move_unknown(ctx, mode)
	local senior, junior = ctx:bot(0), ctx:bot(1)
	local before = prepare(ctx)
	if before == nil then
		return nil
	end
	if fault.rpc(senior, "Family", mode) == false then
		ctx:fail("Family 장애 주입 실패")
		return nil
	end
	send(senior, teleport(senior, junior))
	ctx:sleep(1000)
	fault.rpc(senior, "Family", "off")
	if check(ctx, senior:map() == ELLINIA, "장애 중 이동함: " .. tostring(senior:map())) == false then
		return nil
	end
	return before
end

local function move_lost_timer(ctx)
	local senior = ctx:bot(0)
	local before = move_unknown(ctx, "lost")
	if before == nil then
		return false
	end
	if holds(ctx, senior, before - COST, "장애 직후 (이동 예약됨)") == false then
		return false
	end
	if wait_holds(ctx, senior, before, "주기 정리로 명성 환불") == false then
		return false
	end
	return check(ctx, senior:map() == ELLINIA, "환불 뒤 이동함: " .. tostring(senior:map()))
end

local function move_lost_relogin(ctx)
	local senior = ctx:bot(0)
	local before = move_unknown(ctx, "lost")
	if before == nil then
		return false
	end
	if senior:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_holds(ctx, senior, before, "재접속 정리로 명성 환불") == false then
		return false
	end
	if senior:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	return holds(ctx, senior, before, "복귀 뒤 (중복 환불 없음)")
end

local function move_not_sent(ctx)
	local senior = ctx:bot(0)
	local before = move_unknown(ctx, "on")
	if before == nil then
		return false
	end
	if holds(ctx, senior, before, "장애 직후 (전송 안 됨)") == false then
		return false
	end
	ctx:sleep(SETTLE_WAIT_MS)
	return holds(ctx, senior, before, "주기 정리 뒤")
end

local function late_move_after_refund(ctx)
	local senior, junior = ctx:bot(0), ctx:bot(1)
	local before = prepare(ctx)
	if before == nil then
		return false
	end
	if fault.rpc_delay(senior, "Family", 12000) == false then
		return ctx:fail("Family 지연 주입 실패")
	end
	send(senior, teleport(senior, junior))
	ctx:sleep(300)
	fault.rpc_delay(junior, "Family", 0)

	if senior:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	ctx:sleep(2000)
	if senior:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(14000)
	if check(ctx, senior:map() == ELLINIA, "중단한 뒤 늦은 이동이 반영됨: " .. tostring(senior:map())) == false then
		return false
	end
	return holds(ctx, senior, before, "늦은 요청 뒤")
end

test_suite {
	name = "Family: 순간이동 응답 유실 복구",
	bot_count = 2,

	scenarios = {
		move_lost_timer,
		move_lost_relogin,
		move_not_sent,
		late_move_after_refund,
	},
}

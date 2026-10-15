local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local JAIL = 910000000
local SNAIL = 100100
local DETECTOR = 2190000
local QUIET_MS = 1500

local RESULT = {
	not_found = 0,
	not_hunting = 1,
	tested = 2,
	busy = 3,
	caught = 7,
	passed = 9,
	reported = 11,
}

local seen = {}

local function watch(ctx)
	seen = {}
	ctx:hook("lie_detector", function(_, bot, p, name)
		if name == "lie_detector_started" or name == "lie_detector_result" then
			table.insert(seen, { id = bot:id(), name = name, p = p })
		end
	end)
end

local function received(ctx, bot, name, match)
	for _ = 1, 30 do
		for _, e in ipairs(seen) do
			if e.id == bot:id() and e.name == name and match(e.p) then
				return true
			end
		end
		ctx:sleep(100)
	end
	return false
end

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, meso, items)
	local text = string.format("/봇초기화 30 0 %d %s -", meso, items)
	if pq.command(bot, text, "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:instance_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 헤네시스 이동 실패")
	end
	return true
end

local function meso_is(ctx, bot, meso, what)
	for _ = 1, 20 do
		if bot:meso() == meso then
			return true
		end
		ctx:sleep(100)
	end
	return ctx:fail(string.format("%s %s: 메소 %d (기대 %d)", bot:name(), what, bot:meso(), meso))
end

local function detectors_are(ctx, bot, count, message)
	for _ = 1, 20 do
		if (bot:items()[DETECTOR] or 0) == count then
			return true
		end
		ctx:sleep(100)
	end
	return ctx:fail(message)
end

local function hunt(ctx, bot)
	if bot:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL .. " 1" }) == false then
		return ctx:fail("달팽이 생성 실패")
	end
	local mobs = bot:mobs(SNAIL)
	if check(ctx, #mobs > 0, "생성한 달팽이를 찾지 못함") == false then
		return false
	end
	if check(ctx, bot:attack(mobs[1].oid, 1, 1), "달팽이 공격 실패") == false then
		return false
	end
	ctx:sleep(300)
	return true
end

local function use_item(tester, target)
	return req.lie_detector_item { target = target:name(), slot = tester:slot(DETECTOR) or 0, item_id = DETECTOR }
end

local function result_on(tester, listener, pkt, result, timeout)
	return tester:request_on(listener, resp.lie_detector_result, pkt, function(p)
		return p.result == result
	end, timeout or 3000)
end

local function question(tester, target, pkt)
	return tester:request_on(target, resp.lie_detector_question, pkt, nil, 3000)
end

local function answer_of(tester, target)
	local p = pq.command(tester, "/거탐정답 " .. target:name(), "거탐정답: ")
	if p == false then
		return nil
	end
	return p.message:match("거탐정답: (%w+)")
end

local function rejections(ctx)
	local tester = ctx:bot(0)
	local target = ctx:bot(1)
	if setup(ctx, tester, 0, DETECTOR .. ":3") == false or setup(ctx, target, 0, "-") == false then
		return false
	end

	if check(ctx, result_on(tester, tester, use_item(tester, target), RESULT.not_hunting) ~= false, "사냥하지 않는 대상에게 사냥 중 아님 응답이 없음") == false then
		return false
	end
	if detectors_are(ctx, tester, 3, "거부된 거짓말 탐지기가 소모됨") == false then
		return false
	end
	if check(ctx, result_on(tester, tester, req.lie_detector_item { target = "없는캐릭터", slot = tester:slot(DETECTOR) or 0, item_id = DETECTOR }, RESULT.not_found) ~= false, "없는 대상에게 캐릭터 없음 응답이 없음") == false then
		return false
	end
	if check(ctx, target:request(resp.lie_detector_result, req.lie_detector_answer { answer = "ABCDEF" }, nil, QUIET_MS) == false, "테스트 중이 아닌데 답변이 처리됨") == false then
		return false
	end

	if pq.command(tester, "/플레이어모드", "플레이어 모드: enabled") == false then
		return ctx:fail("플레이어 모드 설정 실패")
	end
	local asked = question(tester, target, req.lie_detector_skill { target = target:name() })
	if pq.command(tester, "/플레이어모드", "플레이어 모드: disabled") == false then
		return ctx:fail("플레이어 모드 해제 실패")
	end
	return check(ctx, asked == false, "일반 유저가 거짓말 탐지 스킬을 사용함")
end

local function pass(ctx)
	local tester = ctx:bot(0)
	local target = ctx:bot(1)
	if setup(ctx, tester, 0, DETECTOR .. ":3") == false or setup(ctx, target, 0, "-") == false then
		return false
	end
	if hunt(ctx, target) == false then
		return false
	end

	watch(ctx)
	local asked = question(tester, target, use_item(tester, target))
	if check(ctx, asked ~= false, "사냥 중인 대상에게 질문이 가지 않음") == false then
		return false
	end
	if check(ctx, asked.skill == false, "아이템 질문이 운영자 창으로 열림") == false then
		return false
	end
	if check(ctx, received(ctx, tester, "lie_detector_started", function(p)
		return p.target == target:name()
	end), "테스터에게 사용 알림이 오지 않음") == false then
		return false
	end
	if detectors_are(ctx, tester, 2, "거짓말 탐지기가 소모되지 않음") == false then
		return false
	end
	if check(ctx, result_on(tester, tester, use_item(tester, target), RESULT.busy) ~= false, "테스트 중인 대상에게 진행 중 응답이 없음") == false then
		return false
	end

	local answer = answer_of(tester, target)
	if check(ctx, answer ~= nil and #answer == 6, "정답 조회 실패") == false then
		return false
	end
	local passed = result_on(target, target, req.lie_detector_answer { answer = answer:lower() }, RESULT.passed)
	if check(ctx, passed ~= false and passed.skill == false, "정답을 맞혔는데 아이템 통과 응답이 없음") == false then
		return false
	end
	if meso_is(ctx, target, 5000, "통과 보상") == false then
		return false
	end
	if check(ctx, answer_of(tester, target) == nil, "통과 후에도 테스트가 남음") == false then
		return false
	end

	if hunt(ctx, target) == false then
		return false
	end
	if check(ctx, result_on(tester, tester, use_item(tester, target), RESULT.tested) ~= false, "통과 직후 대상에게 이미 통과 응답이 없음") == false then
		return false
	end
	return detectors_are(ctx, tester, 2, "거부된 거짓말 탐지기가 소모됨")
end

local function caught(ctx)
	local tester = ctx:bot(0)
	local target = ctx:bot(1)
	if setup(ctx, tester, 0, "-") == false or setup(ctx, target, 10000, "-") == false then
		return false
	end

	watch(ctx)
	local asked = question(tester, target, req.lie_detector_skill { target = target:name() })
	if check(ctx, asked ~= false and asked.skill, "스킬 질문이 운영자 창으로 열리지 않음") == false then
		return false
	end
	if check(ctx, target:request(resp.lie_detector_question, req.lie_detector_answer { answer = "WRONG1" }, nil, 3000) ~= false, "첫 오답 후 새 질문이 오지 않음") == false then
		return false
	end
	local caught = result_on(target, target, req.lie_detector_answer { answer = "WRONG2" }, RESULT.caught)
	if check(ctx, caught ~= false and caught.skill, "두 번 틀렸는데 운영자 단속 적발 응답이 없음") == false then
		return false
	end
	if check(ctx, received(ctx, tester, "lie_detector_result", function(p)
		return p.result == RESULT.reported
	end), "테스터에게 적발 알림이 오지 않음") == false then
		return false
	end
	for _ = 1, 30 do
		if target:map() == JAIL then
			break
		end
		ctx:sleep(100)
	end
	if check(ctx, target:map() == JAIL, "적발된 대상이 이동하지 않음: " .. target:map()) == false then
		return false
	end
	return meso_is(ctx, target, 3000, "적발 차감") and meso_is(ctx, tester, 7000, "적발 보상")
end

local function timeout(ctx)
	local tester = ctx:bot(0)
	local target = ctx:bot(1)
	if setup(ctx, tester, 0, "-") == false or setup(ctx, target, 5000, "-") == false then
		return false
	end
	if pq.command(target, "/타이머 1", "타이머 제한: 1") == false then
		return ctx:fail("타이머 제한 설정 실패")
	end

	local first = question(tester, target, req.lie_detector_skill { target = target:name() })
	local second = first ~= false and target:request(resp.lie_detector_question, nil, nil, 3000)
	local result = second ~= false and target:request(resp.lie_detector_result, nil, function(p)
		return p.result == RESULT.caught
	end, 3000)
	if pq.command(target, "/타이머 0", "타이머 제한: 0") == false then
		return ctx:fail("타이머 제한 해제 실패")
	end
	if check(ctx, first ~= false, "스킬로 질문이 가지 않음") == false then
		return false
	end
	if check(ctx, second ~= false, "첫 시간 초과 후 새 질문이 오지 않음") == false then
		return false
	end
	if check(ctx, result ~= false, "두 번 시간 초과했는데 적발되지 않음") == false then
		return false
	end
	return meso_is(ctx, target, 0, "시간 초과 차감") and meso_is(ctx, tester, 5000, "시간 초과 보상")
end

local function escape(ctx)
	local tester = ctx:bot(0)
	local target = ctx:bot(1)
	if setup(ctx, tester, 0, "-") == false or setup(ctx, target, 20000, "-") == false then
		return false
	end

	if check(ctx, question(tester, target, req.lie_detector_skill { target = target:name() }) ~= false, "스킬로 질문이 가지 않음") == false then
		return false
	end
	if check(ctx, target:transfer(1), "채널 이동 실패") == false then
		return false
	end
	if meso_is(ctx, tester, 7000, "채널 이동 도주 보상") == false or meso_is(ctx, target, 13000, "채널 이동 도주 차감") == false then
		return false
	end
	if check(ctx, target:transfer(0), "채널 복귀 실패") == false then
		return false
	end
	if target:instance_move(HENESYS) == false then
		return ctx:fail("헤네시스 복귀 실패")
	end

	if check(ctx, question(tester, target, req.lie_detector_skill { target = target:name() }) ~= false, "재시험 질문이 가지 않음") == false then
		return false
	end
	if check(ctx, target:relog({ passwords = { "" } }) ~= false, "재접속 실패") == false then
		return false
	end
	return meso_is(ctx, tester, 14000, "접속 종료 도주 보상") and meso_is(ctx, target, 6000, "접속 종료 도주 차감")
end

test_suite {
	name = "LieDetector: 거짓말 탐지기 거부·통과·적발·시간 초과·도주",
	bot_count = 2,

	scenarios = {
		rejections,
		pass,
		caught,
		timeout,
		escape,
	},
}

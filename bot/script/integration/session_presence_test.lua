local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000

-- A lost presence notice must not turn a channel change into a failed login.
local function enter_without_mq(ctx)
	local a, b = ctx:bot(0), ctx:bot(1)
	for _, bot in ipairs({ a, b }) do
		if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
			return ctx:fail(bot:name() .. " 봇 초기화 실패")
		end
		if bot:map_move(HENESYS) == false then
			return ctx:fail(bot:name() .. " 맵 이동 실패")
		end
	end
	if pq.form_party(ctx) == false then
		return false
	end

	if pq.command(a, "/MQ장애 on", "MQ 장애: ON") == false then
		return ctx:fail("MQ 장애 설정 실패")
	end
	local moved = a:transfer(1)
	if pq.command(b, "/MQ장애 off", "MQ 장애: OFF") == false then
		return ctx:fail("MQ 장애 해제 실패")
	end
	if moved == false then
		return ctx:fail("MQ 발행이 끊긴 동안 파티원이 채널을 옮기지 못함")
	end

	if a:transfer(0) == false then
		return ctx:fail("MQ 복구 뒤 원래 채널로 돌아오지 못함")
	end
	for _, bot in ipairs({ b, a }) do
		bot:request(resp.party_update_leave, req.party_operation { operation = PARTY.Leave }, nil, 3000)
	end
	return true
end

test_suite {
	name = "Session presence: MQ 발행 장애 중 채널 이동",
	bot_count = 2,

	scenarios = {
		enter_without_mq,
	},
}

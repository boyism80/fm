local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")

local GUILD_HQ = 200000301
local GUILD_SHOW_INFO = 0x1A
local GUILD_EMBLEM_CHANGE = 0x42

local MESO = 10000000
local CREATE_COST = 1500000
local EMBLEM_COST = 5000000
local SETTLE_WAIT_MS = 45000

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function code(value)
	return function(p)
		return p.code == value
	end
end

local function setup(ctx, bot)
	if pq.command(bot, string.format("/봇초기화 50 0 %d - -", MESO), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= GUILD_HQ and bot:map_move(GUILD_HQ) == false then
		return ctx:fail(bot:name() .. " 길드 본부 이동 실패")
	end
	return true
end

local function holds(ctx, bot, meso, what)
	return check(ctx, bot:meso() == meso, string.format("%s %s: 메소 %d (기대 %d)", bot:name(), what, bot:meso(), meso))
end

local function create_request(bot)
	return req.guild_operation { operation = GUILD.Create, guild_name = bot:name() }
end

local function emblem_request()
	return req.guild_operation { operation = GUILD.ChangeEmblem, logo_bg = 1001, logo_bg_color = 2, logo = 2001, logo_color = 3 }
end

-- The fault is server wide, so it is turned off right after the one request it targets.
local function send_unknown(ctx, bot, method, request, mode)
	if fault.rpc(bot, method, mode) == false then
		return ctx:fail(method .. " 장애 주입 실패")
	end
	bot:send(request)
	ctx:sleep(1500)
	fault.rpc(bot, method, "off")
	return true
end

local function change_emblem(ctx, bot)
	if bot:request(resp.guild_message, emblem_request(), code(GUILD_EMBLEM_CHANGE), 10000) == false then
		return ctx:fail(bot:name() .. " 길드 마크 변경 알림 없음")
	end
	return true
end

local function create_lost_timer(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot) == false then
		return false
	end
	if send_unknown(ctx, bot, "CreateGuild", create_request(bot), "lost") == false then
		return false
	end
	if holds(ctx, bot, MESO - CREATE_COST, "장애 직후 보류") == false then
		return false
	end
	if bot:request(resp.guild_message, nil, code(GUILD_SHOW_INFO), SETTLE_WAIT_MS) == false then
		return ctx:fail("주기 정리 뒤 길드 정보가 오지 않음")
	end
	return holds(ctx, bot, MESO - CREATE_COST, "주기 정리 뒤 (환불 없음)")
end

local function create_lost_relogin(ctx)
	local bot = ctx:bot(1)
	if setup(ctx, bot) == false then
		return false
	end
	if send_unknown(ctx, bot, "CreateGuild", create_request(bot), "lost") == false then
		return false
	end
	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	ctx:sleep(3000)
	if bot:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, bot, MESO - CREATE_COST, "재접속 정리 뒤 (환불 없음)") == false then
		return false
	end
	if change_emblem(ctx, bot) == false then
		return false
	end
	ctx:sleep(500)
	return holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "길드 마크 변경 뒤")
end

local function create_not_sent(ctx)
	local bot = ctx:bot(2)
	if setup(ctx, bot) == false then
		return false
	end
	if send_unknown(ctx, bot, "CreateGuild", create_request(bot), "on") == false then
		return false
	end
	if bot:request(resp.guild_message, nil, code(GUILD_SHOW_INFO), SETTLE_WAIT_MS) == false then
		return ctx:fail("주기 정리로 다시 보낸 길드 생성 정보가 오지 않음")
	end
	return holds(ctx, bot, MESO - CREATE_COST, "다시 보낸 생성 뒤")
end

local function late_create_after_resend(ctx)
	local bot = ctx:bot(3)
	local other = ctx:bot(4)
	if setup(ctx, bot) == false then
		return false
	end
	if fault.rpc_delay(bot, "CreateGuild", 6000) == false then
		return ctx:fail("CreateGuild 지연 주입 실패")
	end
	bot:send(create_request(bot))
	local moved = bot:transfer(1)
	fault.rpc_delay(other, "CreateGuild", 0)
	if moved == false then
		return ctx:fail("채널 이동 실패")
	end

	ctx:sleep(10000)
	if holds(ctx, bot, MESO - CREATE_COST, "재접속 재요청 뒤") == false then
		return false
	end
	ctx:sleep(8000)
	if holds(ctx, bot, MESO - CREATE_COST, "늦은 원래 요청 뒤 (환불·중복 차감 없음)") == false then
		return false
	end
	if bot:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	return change_emblem(ctx, bot)
end

local function emblem_lost(ctx)
	local bot = ctx:bot(0)
	if send_unknown(ctx, bot, "ChangeGuildEmblem", emblem_request(), "lost") == false then
		return false
	end
	if holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "장애 직후 보류") == false then
		return false
	end
	ctx:sleep(SETTLE_WAIT_MS)
	return holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "주기 정리 뒤 (환불 없음)")
end

local function emblem_not_sent(ctx)
	local bot = ctx:bot(2)
	if send_unknown(ctx, bot, "ChangeGuildEmblem", emblem_request(), "on") == false then
		return false
	end
	if bot:request(resp.guild_message, nil, code(GUILD_EMBLEM_CHANGE), SETTLE_WAIT_MS) == false then
		return ctx:fail("주기 정리로 다시 보낸 길드 마크 변경 알림 없음")
	end
	ctx:sleep(500)
	return holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "다시 보낸 마크 변경 뒤")
end

test_suite {
	name = "Guild: 길드 생성·마크 변경 응답 유실 복구",
	bot_count = 5,

	scenarios = {
		create_lost_timer,
		create_lost_relogin,
		create_not_sent,
		late_create_after_resend,
		emblem_lost,
		emblem_not_sent,
	},
}

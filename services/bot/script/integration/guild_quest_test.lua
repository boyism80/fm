local pq = require("script/integration/lib/party_quest")

local GUILD_HQ = 200000301
local CAMP = 101030104
local EXCAVATION = 990000000
local SHUANG = 9040000

local GUILD_SHOW_INFO = 0x1A
local DIALOG_LIST = 4

local function show_info(p)
	return p.code == GUILD_SHOW_INFO
end

local function register(ctx, bot, selected)
	local oid = pq.npc(ctx, bot, SHUANG)
	if oid == false then
		return false
	end
	if bot:npc_click(oid) == false then
		return ctx:fail(bot:name() .. " 슈앵 대화가 오지 않음")
	end
	if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = selected }, function(p)
			return p.character.map == EXCAVATION
		end, 10000) == false then
		return ctx:fail(bot:name() .. " 유적발굴 현장 입장 실패")
	end
	return true
end

test_suite {
	name = "Guild: 길드 대항전 접수",
	bot_count = 6,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/봇초기화 50 100 2000000", "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(GUILD_HQ) == false then
				return ctx:fail(bot:name() .. " 길드 본부 이동 실패")
			end
		end
		return true
	end,

	scenarios = {
		function(ctx)
			local leader = ctx:bot(0)
			if leader:request(resp.guild_message, req.guild_operation { operation = GUILD.Create, guild_name = leader:name() }, show_info) == false then
				return ctx:fail("길드 생성 실패")
			end
			for i = 1, ctx:bot_count() - 1 do
				local member = ctx:bot(i)
				local invite = leader:request_on(member, resp.guild_invite,
					req.guild_operation { operation = GUILD.Invite, target_name = member:name() })
				if invite == false then
					return ctx:fail(member:name() .. " 길드 초대 패킷 없음")
				end
				local joined = member:request(resp.guild_message, req.guild_operation {
					operation = GUILD.AcceptInvite,
					guild_id = invite.guild_id,
					character_id = member:id(),
				}, show_info)
				if joined == false then
					return ctx:fail(member:name() .. " 길드 가입 실패")
				end
			end
			return true
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				if ctx:bot(i):map_move(CAMP) == false then
					return ctx:fail(ctx:bot(i):name() .. " 유적발굴단 캠프 이동 실패")
				end
			end
			if register(ctx, ctx:bot(0), 0) == false then
				return false
			end
			for i = 1, ctx:bot_count() - 1 do
				if register(ctx, ctx:bot(i), 1) == false then
					return false
				end
			end
			return true
		end,
		function(ctx)
			local opened = ctx:bot(0):request(resp.notice, nil, function(p)
				return p.message:find("샤레니안의 문이 열렸습니다", 1, true) ~= nil
			end, 200000)
			if opened == false then
				return ctx:fail("3분 뒤 샤레니안의 문이 열리지 않음")
			end
			return true
		end,
		function(ctx)
			for i = ctx:bot_count() - 1, 0, -1 do
				if ctx:bot(i):map_move(CAMP) == false then
					return ctx:fail(ctx:bot(i):name() .. " 캠프로 돌아가기 실패")
				end
			end
			return true
		end,
	},
}

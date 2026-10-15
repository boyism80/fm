local pq = require("script/integration/lib/party_quest")

local RUSH_ENTRY = 105090700
local RUSH_FIELD = 910500100
local BERSERK_ENTRY = 105090800
local BERSERK_FIELD = 910500200
local BERSERK_KEY = 4031475

test_suite {
	name = "Fourth Job: 전사 러시·버서크 시험",
	bot_count = 2,

	on_initialize = function(ctx)
		local classes = { 112, 122 }
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local command = string.format("/봇초기화 120 %d 0 - 6110:1", classes[i + 1])
			if pq.command(bot, command, "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(RUSH_ENTRY) == false then
				return ctx:fail(bot:name() .. " 러시 입장 맵 이동 실패")
			end
		end
		return true
	end,

	on_finished = function(ctx)
		ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	end,

	scenarios = {
		function(ctx)
			local leader = ctx:bot(0)
			leader:warp("in00")
			if leader:map() ~= RUSH_ENTRY then
				return ctx:fail("파티 없이 러시 시험장에 입장함")
			end
			return true
		end,
		pq.form_party,
		function(ctx)
			local leader = ctx:bot(0)
			if pq.portal(ctx, leader, "in00", RUSH_FIELD) == false then
				return false
			end
			return pq.wait_all(ctx, RUSH_FIELD, 10000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local mob = nil
			for _ = 1, 10 do
				mob = leader:mobs(9300087)[1]
				if mob ~= nil then
					break
				end
				ctx:sleep(200)
			end
			if mob == nil then
				return ctx:fail("러시 시험 몬스터가 보이지 않음")
			end
			return leader:kill(mob.oid)
		end,
		function(ctx)
			local member = ctx:bot(1)
			if member:request(resp.party_update_leave, req.party_operation { operation = PARTY.Leave }, nil, 5000) == false then
				return ctx:fail("러시 시험 중 파티 탈퇴 실패")
			end
			return pq.wait_all(ctx, RUSH_ENTRY, 10000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			bot:request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)

			local command = string.format("/봇초기화 120 112 0 %d:1 6153:1", BERSERK_KEY)
			if pq.command(bot, command, "봇초기화 완료") == false then
				return ctx:fail("버서크 시험 봇 초기화 실패")
			end
			if bot:map_move(BERSERK_ENTRY) == false then
				return ctx:fail("버서크 입장 맵 이동 실패")
			end
			if pq.portal(ctx, bot, "in00", BERSERK_FIELD) == false then
				return false
			end
			if (bot:items()[BERSERK_KEY] or 0) ~= 1 then
				return ctx:fail("버서크 입장 열쇠가 소비됨")
			end
			for _ = 1, 10 do
				local mob = bot:mobs(9300092)[1]
				if mob ~= nil then
					return bot:kill(mob.oid)
				end
				ctx:sleep(200)
			end
			return ctx:fail("버서크 시험 몬스터가 보이지 않음")
		end,
	},
}

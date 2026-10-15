local pq = require("script/integration/lib/party_quest")

local NAVIGATION_ROOM = 120000101
local TEST_MAPS = { [108000502] = true, [108000503] = true }
local KYRIN = 1090000
local KYRIN_INSIDE = 1072008
local OCTO = 9001006
local PROOF = 4031856

local DIALOG_YES_NO = 1

test_suite {
	name = "Job: 카이린의 2차 전직 시험",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 30 500 0 - 2191:1", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(NAVIGATION_ROOM)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, KYRIN)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("카이린 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
					return TEST_MAPS[p.character.map] ~= nil
				end, 10000) == false then
				return ctx:fail("해적의 시험장 입장 실패")
			end
			return true
		end,
		function(ctx)
			return pq.collect(ctx, ctx:bot(0), OCTO, PROOF, 15)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, KYRIN_INSIDE)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("시험장 카이린 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
					return p.character.map == NAVIGATION_ROOM
				end, 10000) == false then
				return ctx:fail("항해실로 돌아가지 않음")
			end
			return true
		end,
	},
}

local pq = require("script/integration/lib/party_quest")

local EOS_100 = 221024400
local DOLL_HOUSE = 922000010
local OLSON = 2040002
local MARK = 2040028
local ODD_HOUSE = 2202000
local PENDULUM = 4031145

local DIALOG_DEFAULT = 0

test_suite {
	name = "Quest: 인형의 집",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 50 100 0 - 3230:1", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(EOS_100)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, OLSON)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("올슨 대화가 오지 않음")
			end
			if bot:dialog(true) == false then
				return ctx:fail("올슨 안내 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == DOLL_HOUSE
				end, 10000) == false then
				return ctx:fail("인형의 집 입장 실패")
			end
			return true
		end,
		function(ctx)
			return pq.open_box(ctx, ctx:bot(0), ODD_HOUSE, PENDULUM)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, MARK)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("마크 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == EOS_100
				end, 10000) == false then
				return ctx:fail("시계추를 건넨 뒤 에오스탑 100층으로 나가지 않음")
			end
			if (bot:items()[PENDULUM] or 0) > 0 then
				return ctx:fail("시계추가 회수되지 않음")
			end
			return true
		end,
	},
}

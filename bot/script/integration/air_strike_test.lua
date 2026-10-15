local pq = require("script/integration/lib/party_quest")

local JONATHAN_ROOM = 120000102
local TRAINING_ROOM = 120000104
local BART_ROOM = 912020000
local GATE = 910210000
local JONATHAN = 1092019
local STATUE = 1022103
local REAL_BART = 1209000

local DIALOG_DEFAULT = 0
local DIALOG_YES_NO = 1

local function repeat_flashes(ctx, bot, count)
	local order = {}
	for i = 1, count do
		local p = bot:request(resp.trigger_reactor, nil, nil, 8000)
		if p == false then
			return ctx:fail(string.format("동상이 %d번째로 반짝이지 않음", i))
		end
		order[i] = p.reactor.oid
	end
	for _, oid in ipairs(order) do
		local reactor = pq.find_reactor(bot, oid)
		if reactor == nil then
			return ctx:fail("반짝인 동상이 보이지 않음: " .. oid)
		end
		if pq.move(bot, reactor.x, reactor.y) == false then
			return ctx:fail("동상 앞으로 이동 실패: " .. reactor.name)
		end
		if bot:hit_reactor(oid) == false then
			return ctx:fail("동상 타격 응답 없음: " .. reactor.name)
		end
	end
	return true
end

test_suite {
	name = "Quest: 조나단의 시험 (에어 스트라이크)",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 70 500 0 - 6400:1 - air_strike.progress=q2", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(JONATHAN_ROOM)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, JONATHAN)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("조나단 대화가 오지 않음")
			end
			if bot:dialog(true) == false then
				return ctx:fail("조나단 준비 확인이 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
					return p.character.map == BART_ROOM
				end, 10000) == false then
				return ctx:fail("바트의 방 입장 실패")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local bart = pq.seek_reactor(ctx, bot, pq.reactor_by_id(REAL_BART))
			if bart == false then
				return false
			end
			if pq.break_reactor(ctx, bot, bart) == false then
				return false
			end
			return pq.wait_map(ctx, bot, TRAINING_ROOM, 10000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if bot:map_move(JONATHAN_ROOM) == false then
				return ctx:fail("조나단의 방 이동 실패")
			end
			local oid = pq.npc(ctx, bot, JONATHAN)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("조나단 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == GATE
				end, 10000) == false then
				return ctx:fail("샤레니안 성문 입장 실패")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, STATUE)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("분수 조각상 대화가 오지 않음")
			end
			if bot:dialog(true) == false then
				return ctx:fail("분수 조각상 안내가 오지 않음")
			end
			bot:dialog(false)
			for stage = 1, 3 do
				local dlg = bot:npc_click(oid)
				if dlg == false then
					return ctx:fail(stage .. "단계 대화가 오지 않음")
				end
				bot:dialog(false)
				if repeat_flashes(ctx, bot, stage + 3) == false then
					return false
				end
				if pq.move(bot, 621, 156) == false then
					return ctx:fail("분수 조각상 앞으로 이동 실패")
				end
				if stage < 3 then
					local result = bot:npc_click(oid)
					if result == false then
						return ctx:fail(stage .. "단계 결과 대화가 오지 않음")
					end
					bot:dialog(false)
					if result.text:find("올바르게", 1, true) == nil then
						return ctx:fail(stage .. "단계 실패: " .. result.text)
					end
				end
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("마지막 결과 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == JONATHAN_ROOM
				end, 10000) == false then
				return ctx:fail("시험을 통과한 뒤 조나단의 방으로 돌아가지 않음")
			end
			return true
		end,
	},
}

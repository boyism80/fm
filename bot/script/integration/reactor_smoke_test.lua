local pq = require("script/integration/lib/party_quest")

local MAPLE_ISLAND = 1000000
local BOX = 2000
local FOLK_TOWN_FIELD = 222010401
local ROCK_PILE = 2221000
local SKEWER = 2022050
local YELLOW_GOBLIN = 7130400

local function approach_reactor(ctx, bot, template)
	for _, spot in ipairs(bot:reactor_spots()) do
		if spot.id == template then
			if pq.move(bot, spot.x, spot.y) == false then
				return ctx:fail("리액터 앞으로 이동 실패: " .. template)
			end
			local reactor = bot:reactors(template)[1]
			if reactor == nil then
				return ctx:fail("리액터가 보이지 않음: " .. template)
			end
			return reactor
		end
	end
	return ctx:fail("맵에 리액터 배치가 없음: " .. template)
end

test_suite {
	name = "Reactor Smoke (타격, 아이템 리액터)",
	bot_count = 1,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if bot:instance_move(MAPLE_ISLAND) == false then
				return ctx:fail("메이플 아일랜드 이동 실패")
			end
			local reactor = approach_reactor(ctx, bot, BOX)
			if reactor == false then
				return false
			end

			local state = reactor.state
			for _ = 1, 10 do
				local p, name = bot:hit_reactor(reactor.oid)
				if p == false then
					return ctx:fail("리액터 타격 응답 없음")
				end
				if name == resp.destroy_reactor then
					local drop = bot:request({ resp.spawn_item, resp.spawn_meso }, nil, nil, 3000)
					if drop == false then
						return ctx:fail("상자가 아무것도 떨어뜨리지 않음")
					end
					return true
				end
				if p.reactor.state <= state then
					return ctx:fail("리액터 상태가 진행되지 않음: " .. p.reactor.state)
				end
				state = p.reactor.state
			end
			return ctx:fail("10번 때려도 상자가 부서지지 않음")
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if bot:instance_move(FOLK_TOWN_FIELD) == false then
				return ctx:fail("한국민속촌 필드 이동 실패")
			end
			if pq.command(bot, "/아이템생성 " .. SKEWER .. " 1", "아이템 생성") == false then
				return ctx:fail("돼지고기 산적 생성 실패")
			end
			local reactor = approach_reactor(ctx, bot, ROCK_PILE)
			if reactor == false then
				return false
			end

			if bot:drop(SKEWER, 1) == nil then
				return ctx:fail("돼지고기 산적 버리기 실패")
			end
			local p = bot:request({ resp.trigger_reactor, resp.destroy_reactor }, nil, function(p)
				return p.reactor.oid == reactor.oid
			end, 10000)
			if p == false then
				return ctx:fail("돌무더기가 반응하지 않음")
			end
			if #bot:mobs(YELLOW_GOBLIN) == 0 then
				local spawn = bot:request(resp.spawn_mob, nil, function(p)
					return p.mob.mob_id == YELLOW_GOBLIN
				end, 5000)
				if spawn == false then
					return ctx:fail("노란왕도깨비가 소환되지 않음")
				end
			end
			return true
		end,
	},
}

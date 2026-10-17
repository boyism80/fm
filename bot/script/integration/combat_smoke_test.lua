local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local SNAIL = 100100
local SNAIL_SHELL = 4000019
local TIMED_MOB = 9300117
local TIMED_MOB_SECONDS = 10
local PART_HP = 700
local SWORD = 1302000
local WEAPON_SLOT = -11

local function command(bot, name, text)
	return bot:request(name, req.normal_chat { message = text })
end

test_suite {
	name = "Combat Smoke (몹 처치, 드롭, 줍기, 공격자 없는 처치)",
	bot_count = 1,

	on_finished = function(ctx)
		pq.command(ctx:bot(0), "/즉사 on", "즉사 상태: enabled")
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if bot:instance_move(HUNTING_GROUND) == false then
				return ctx:fail("사냥터 이동 실패")
			end

			local before = {}
			for _, mob in ipairs(bot:mobs(SNAIL)) do
				before[mob.oid] = true
			end
			local spawn = command(bot, resp.spawn_mob, "/몬스터생성 " .. SNAIL)
			if spawn == false then
				return ctx:fail("달팽이 생성 응답 없음")
			end
			local oid = nil
			for _, mob in ipairs(bot:mobs(SNAIL)) do
				if before[mob.oid] == nil then
					oid = mob.oid
				end
			end
			if oid == nil then
				return ctx:fail("생성한 달팽이를 찾지 못함")
			end
			if bot:kill(oid) == false then
				return ctx:fail("달팽이 처치 실패")
			end
			for _, mob in ipairs(bot:mobs(SNAIL)) do
				if mob.oid == oid then
					return ctx:fail("처치한 달팽이가 남아 있음")
				end
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if command(bot, resp.inventory_operation, "/아이템생성 " .. SNAIL_SHELL .. " 5") == false then
				return ctx:fail("아이템 생성 응답 없음")
			end
			local count = bot:items()[SNAIL_SHELL] or 0
			if count < 5 then
				return ctx:fail("아이템 생성 후 개수: " .. count)
			end

			local oid = bot:drop(SNAIL_SHELL, 2)
			if oid == nil then
				return ctx:fail("아이템 드롭 실패")
			end
			if (bot:items()[SNAIL_SHELL] or 0) ~= count - 2 then
				return ctx:fail("드롭 후 개수: " .. (bot:items()[SNAIL_SHELL] or 0))
			end

			if bot:loot(oid) == false then
				return ctx:fail("아이템 줍기 실패")
			end
			for _, drop in ipairs(bot:drops(SNAIL_SHELL)) do
				if drop.oid == oid then
					return ctx:fail("주운 아이템이 맵에 남아 있음")
				end
			end
			if (bot:items()[SNAIL_SHELL] or 0) ~= count then
				return ctx:fail("줍기 후 개수: " .. (bot:items()[SNAIL_SHELL] or 0))
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if pq.command(bot, "/봇초기화 62 110 0 " .. SWORD .. ":1 -", "봇초기화 완료") == false then
				return ctx:fail("봇 초기화 실패")
			end
			if bot:request(resp.inventory_operation, req.move_item { inventory_type = 1, source = bot:slot(SWORD) or 0, dest = WEAPON_SLOT, count = 1 }, nil, 3000) == false then
				return ctx:fail("검 착용 실패")
			end
			if pq.command(bot, "/전능 999", "전능 설정") == false then
				return ctx:fail("스탯 설정 실패")
			end
			if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
				return ctx:fail("즉사 끄기 실패")
			end
			if bot:instance_move(HUNTING_GROUND) == false then
				return ctx:fail("사냥터 이동 실패")
			end
			local before = {}
			for _, mob in ipairs(bot:mobs(TIMED_MOB)) do
				before[mob.oid] = true
			end
			if command(bot, resp.spawn_mob, "/몬스터생성 " .. TIMED_MOB) == false then
				return ctx:fail("몬스터 생성 응답 없음")
			end
			local oid = nil
			for _, mob in ipairs(bot:mobs(TIMED_MOB)) do
				if before[mob.oid] == nil then
					oid = mob.oid
				end
			end
			if oid == nil then
				return ctx:fail("생성한 몬스터를 찾지 못함")
			end
			bot:attack(oid, PART_HP, 1)
			ctx:sleep(500)
			local alive = false
			for _, mob in ipairs(bot:mobs(TIMED_MOB)) do
				if mob.oid == oid then
					alive = true
				end
			end
			if not alive then
				return ctx:fail("일부 피해로 몬스터가 죽음")
			end

			local exp = bot:exp()
			local drops = #bot:drops()
			if bot:request(resp.die_mob, nil, function(p)
					return p.oid == oid
				end, (TIMED_MOB_SECONDS + 5) * 1000) == false then
				return ctx:fail("시간이 지나도 몬스터가 사라지지 않음")
			end
			ctx:sleep(1000)
			if bot:exp() ~= exp then
				return ctx:fail("시간이 지나 사라진 몬스터가 경험치를 줌: " .. (bot:exp() - exp))
			end
			if #bot:drops() ~= drops then
				return ctx:fail("시간이 지나 사라진 몬스터가 드롭함")
			end
			return true
		end,
	},
}

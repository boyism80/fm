local HUNTING_GROUND = 104040000
local SNAIL = 100100
local SNAIL_SHELL = 4000019

local function command(bot, name, text)
	return bot:request(name, req.normal_chat { message = text })
end

test_suite {
	name = "Combat Smoke (몹 처치, 드롭, 줍기)",
	bot_count = 1,

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
	},
}

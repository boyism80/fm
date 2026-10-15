local HENESYS = 100000000

local function move(i)
	return function(ctx)
		if ctx:bot(i):instance_move(HENESYS) == false then
			return ctx:fail("헤네시스 이동 실패: 봇 " .. i)
		end
		return true
	end
end

test_suite {
	name = "Login smoke",
	bot_count = 2,

	scenarios = {
		{ parallel = { move(0), move(1) } },
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				if bot:map() ~= HENESYS then
					return ctx:fail(bot:name() .. " 맵 " .. bot:map())
				end
			end
			return true
		end,
	},
}

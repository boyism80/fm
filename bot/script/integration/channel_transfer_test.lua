local CHANNEL = 1
local HOME = 0

test_suite {
	name = "Channel transfer",
	bot_count = 1,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local map = bot:map()
			if bot:transfer(CHANNEL) == false then
				return ctx:fail("채널 " .. CHANNEL .. " 이동 실패")
			end
			if bot:map() ~= map then
				return ctx:fail("채널 이동 뒤 맵이 바뀜: " .. map .. " -> " .. bot:map())
			end
			if bot:hp() <= 0 then
				return ctx:fail("채널 이동 뒤 HP " .. bot:hp())
			end

			local x, y = bot:position()
			if x == nil then
				return ctx:fail("맵 " .. bot:map() .. "의 스폰 위치가 WZ에 없음")
			end
			log("info", string.format("%s at %d (%d, %d)", bot:name(), bot:map(), x, y))

			if bot:transfer(HOME) == false then
				return ctx:fail("채널 " .. HOME .. " 복귀 실패")
			end
			return true
		end,
	},
}

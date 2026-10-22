-- Map script (onFirstUserEnter StageMsg_together): 단계 안내 문구

local HINTS = {
	[103000800] = "문제를 해결하고 해답에 맞는 개수만큼 쿠폰을 모아오세요!",
	[103000801] = "줄에 매달리고, 올바른 줄을 찾으세요!",
	[103000802] = "발판위에 올라서서, 올바른 발판을 찾으세요!",
	[103000803] = "통 위에 올라서서, 올바른 통을 찾으세요!",
	[103000804] = "킹슬라임과 부하들을 물리치고 통행증을 모아오세요!",
}

return {
	on_activate = function(me, map)
		local hint = HINTS[map:wz():id()]
		if hint ~= nil then
			map:start_weather(5120017, hint)
		end
	end,
}

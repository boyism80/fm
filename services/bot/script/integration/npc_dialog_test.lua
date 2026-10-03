local LAKELIS = 9020000
local KERNING = 103000000

test_suite {
	name = "NPC Dialog (라케리스, 파티 없음)",
	bot_count = 1,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if bot:instance_move(KERNING) == false then
				return ctx:fail("커닝시티 이동 실패")
			end

			local oid = bot:npc(LAKELIS)
			if oid == nil then
				return ctx:fail("라케리스 NPC가 맵에 없음")
			end

			local dlg = bot:npc_click(oid)
			if dlg == false then
				return ctx:fail("대화 응답 없음")
			end
			if dlg.text:find("파티장", 1, true) == nil then
				return ctx:fail("예상과 다른 대화: " .. dlg.text)
			end

			bot:dialog(false)
			return true
		end,
	},
}

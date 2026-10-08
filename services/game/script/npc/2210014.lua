-- NPC name (String.wz/Npc.img.xml): 구와르

local SCRIPT = "script/npc/2210014.lua"
local DEST = 240093300
local DAILY_LIMIT = 3
local MIN_LEVEL = 10

return {
	reset_map = function(map)
		map:reset()
	end,

	on_click = function(me, npc)
		local count = me:records():get("entry.guwar")
		if count >= DAILY_LIMIT then
			me:dialog(npc, "더이상 입장 하실 수 없습니다.")
			return
		end
		if me:party() == nil then
			me:dialog(npc, "파티를 만들어 주세요.")
			return
		end
		local sel = me:dialog_list(npc, "이곳은 정말 위험하다..\r\n", {
			"#b암벽거인의 심장에 입장하겠습니다. (하루 1회)",
		})
		if sel == nil then
			return
		end

		local party = me:party()
		if party == nil then
			return
		end
		for _, mem in pairs(party:members()) do
			if mem:level() < MIN_LEVEL then
				me:dialog(npc, "레벨이 낮거나 횟수가 다한 유저가 있습니다.")
				return
			end
		end
		me:records():add("entry.guwar", 1, { daily = true })
		run_on_map(DEST, SCRIPT, "reset_map")
		local pid = party:id()
		for _, ch in pairs(me:map():characters()) do
			local p = ch:party()
			if p ~= nil and p:id() == pid then
				ch:map(DEST, 0)
			end
		end
	end
}

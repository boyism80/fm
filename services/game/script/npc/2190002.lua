-- NPC name (String.wz/Npc.img.xml): 핫세

local SCRIPT = "script/npc/2190002.lua"
local DEST = 223030210
local DAILY_LIMIT = 1
local MIN_LEVEL = 10

return {
	prepare_stage = function(map)
		map:reset()
		map:spawn_mob(8145100, -25, -4)
	end,

	on_click = function(me, npc)
		local count = me:records():get("entry.hasse")
		if count >= DAILY_LIMIT then
			me:dialog(npc, "더이상 입장 하실 수 없습니다.")
			return
		end
		if me:party() == nil then
			me:dialog(npc, "파티를 만들어 주세요.")
			return
		end
		local sel = me:dialog_list(npc, "환상의 나라 판타스틱 테마파크로~\r\n", {
			"#b공연장에 입장하겠습니다. (하루 1회)",
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
		me:records():add("entry.hasse", 1, { daily = true })
		run_on_map(DEST, SCRIPT, "prepare_stage")
		local pid = party:id()
		for _, ch in pairs(me:map():characters()) do
			local p = ch:party()
			if p ~= nil and p:id() == pid then
				ch:map(DEST, 0)
			end
		end
	end
}

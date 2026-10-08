-- NPC name (String.wz/Npc.img.xml): 유적정령

local MAIN = {
	{ text = "10개#k", npc = 9110002 },
	{ text = "100개#k", npc = 9110009 },
	{ text = "1000개#k", npc = 9075000 },
}

local function choose(me, npc, prompt, entries)
	local options = {}
	for i, entry in ipairs(entries) do
		options[i] = entry.text
	end
	local sel = me:dialog_list(npc, prompt, options)
	if sel == nil then
		return
	end
	local entry = entries[sel]
	if entry.npc ~= nil then
		me:open_npc(entry.npc)
		return
	end
	if id2map(entry.map) == nil then
		me:dialog(npc, "아직 갈 수 없는 곳입니다.")
		return
	end
	me:map(entry.map)
end

return {
	on_click = function(me, npc)
		choose(me, npc, "#k유물 교환 기능 관련입니다.", MAIN)
	end
}

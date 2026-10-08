-- NPC name (String.wz/Npc.img.xml): 가가

local MAIN = {
	{ text = "#b가가와 대화 한다.#k", npc = 9000021 },
	{ text = "#b유물 포인트를 교환 한다.#k", npc = 1540848 },
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
	if entry.menu ~= nil then
		choose(me, npc, "", entry.menu)
		return
	end
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
		choose(me, npc, "#k어느걸 이용할거야?", MAIN)
	end
}

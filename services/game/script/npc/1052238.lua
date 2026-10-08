-- NPC name (String.wz/Npc.img.xml): 사피

local MAIN = {
	{ text = "앱솔 랩스#k", npc = 9000175 },
	{ text = "여제 무기", npc = 1530574 },
	{ text = "여제 방어#k", npc = 2201001 },
	{ text = "사자 무기#k", npc = 2161002 },
	{ text = "발록 무기#k", npc = 1061016 },
	{ text = "메플 무기#k", npc = 2019 },
	{ text = "벚꽃 무기#k", npc = 9110101 },
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
		choose(me, npc, "#k어느걸 이용할거야?", MAIN)
	end
}

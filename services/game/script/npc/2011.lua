-- NPC name (String.wz/Npc.img.xml): 할로 캣

local MAIN = {
	{ text = "티어 이용#r//승급,퀘스트 받기#k", npc = 2009 },
	{ text = "위치 이동#r//사냥터,보스 이동#k", npc = 2010 },
	{ text = "상점 이용#r//부산물 아이템 교환#k", npc = 2012 },
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
		choose(me, npc, "#k반복 퀘스트 진행이 어려우시다구요? #e@이동#n 을 입력후 #e강화소 이동#n 을 이용하여 티어 엔피시 에게 받아주시면 더욱 편한 시스템 이용이 가능합니다.", MAIN)
	end
}

-- NPC name (String.wz/Npc.img.xml): 할로캣

local HUNTS = {
	{ text = "#r[추천 Lv 10~30]#k 크리스 마스 공터 #b(이벤트)#k", map = 209000200 },
	{ text = "#r[추천 Lv 48~80]#k 엘린 숲 #b(트리로드)#k", map = 300020000 },
	{ text = "#r[추천 Lv 80~105]#k 버려진 사유지 1 #b(버서키, 비트론)#k", map = 541020100 },
	{ text = "#r[추천 Lv 105~200]#k 버려진 도시의 중심지 #b(페트리 파이터)#k", map = 541020500 },
}

local PARTIES = {
	{ text = "#r[파티 3~6인]#k 레벨제한 : 10~250 #b월묘 파퀘#k", map = 100000200 },
	{ text = "#r[파티 3~6인]#k 레벨제한 : 21~250 #b커닝 파퀘#k", map = 103000000 },
	{ text = "#r[파티 2~6인]#k 레벨제한 : 30~250 #b몬스터 카니발#k", npc = 2042002 },
	{ text = "#r[파티 5~6인]#k 레벨제한 : 35~250 #b루디 파퀘#k", map = 221024500 },
	{ text = "#r[파티 6~6인]#k 레벨제한 : 51~250 #b올비 파퀘#k", map = 200080101 },
	{ text = "#r[파티 3~6인]#k 레벨제한 : 45~250 #b엘린숲 파퀘#k", map = 300030100 },
	{ text = "#r[파티 3~6인]#k 레벨제한 : 55~250 #b데비존 파퀘#k", map = 251010404 },
	{ text = "#r[파티 4~6인]#k 레벨제한 : 71~250 #b줄리엣 파퀘#k", map = 261000021 },
	{ text = "#r[파티 4~6인]#k 레벨제한 : 71~250 #b로미오 파퀘#k", map = 261000011 },
}

local MAIN = {
	{ text = "#b쩔맵 이동#k", menu = HUNTS },
	{ text = "#b파퀘 이동#k", menu = PARTIES },
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

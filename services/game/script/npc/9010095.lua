-- NPC name (String.wz/Npc.img.xml): 할로캣

local TOWNS = {
	{ text = "자유시장", map = 910000000 },
	{ text = "코크타운", map = 219000000 },
	{ text = "행복한 마을", map = 209000000 },
	{ text = "웨딩빌리지", map = 680000000 },
	{ text = "리스항구", map = 104000000 },
	{ text = "헤네시스", map = 100000000 },
	{ text = "엘리니아", map = 101000000 },
	{ text = "커닝시티", map = 103000000 },
	{ text = "페리온", map = 102000000 },
	{ text = "슬리피우드", map = 105040300 },
	{ text = "노틸러스", map = 120000000 },
	{ text = "플로리나 비치", map = 110000000 },
	{ text = "오르비스", map = 200000000 },
	{ text = "엘나스", map = 211000000 },
	{ text = "루디브리엄", map = 220000000 },
	{ text = "지구방위본부", map = 221000000 },
	{ text = "아랫마을", map = 222000000 },
	{ text = "아쿠아리움", map = 230000000 },
	{ text = "리프레", map = 240000000 },
	{ text = "무릉도원", map = 250000000 },
	{ text = "백초마을", map = 251000000 },
	{ text = "아리안트", map = 260000000 },
	{ text = "마가티아", map = 261000000 },
	{ text = "시간의 신전", map = 270000000 },
	{ text = "엘린숲", map = 300000000 },
	{ text = "태국 - 플로팅마켓", map = 500000000 },
	{ text = "대만 - 서문정", map = 740000000 },
	{ text = "일본 - 버섯신사", map = 800000000 },
	{ text = "중국 - 상해와이탄", map = 701000000 },
	{ text = "오르비스 - 길드본부<영웅의전당>", map = 200000301 },
	{ text = "히든스트리트 - 마르의숲", map = 101000200 },
}

local HUNTS = {
	{ text = "#r[추천 Lv 10~30]#k 크리스 마스 공터 #b(이벤트)#k", map = 209000200 },
	{ text = "#r[추천 Lv 48~80]#k 엘린 숲 #b(트리로드)#k", map = 300020000 },
	{ text = "#r[추천 Lv 80~105]#k 버려진 사유지 1 #b(버서키, 비트론)#k", map = 541020100 },
	{ text = "#r[추천 Lv 105~200]#k 버려진 도시의 중심지 #b(페트리 파이터)#k", map = 541020500 },
}

local BOSSES = {
	{ text = "#r[파풀라투스]#k #b시계탑 깊은 곳#k (소환할때)", map = 220080000 },
	{ text = "#r[파풀라투스]#k #b기계탑의 근원#k #e(팅길때)#n", map = 220080001 },
	{ text = "#r[자쿰]#k #b자쿰의 제단 입구#k (소환할때)", map = 211042400 },
	{ text = "#r[자쿰]#k #b자쿰의 제단#k #e(팅길때)#n", map = 280030000 },
	{ text = "#r[혼테일]#k #b혼테일의 동굴 입구#k (소환할때)", map = 240050400 },
	{ text = "#r[혼테일]#k #b혼테일의 동굴 입구#k #e(팅길때)#n", map = 240060200 },
	{ text = "#r[핑크빈]#k #b잊혀진 황혼#k (소환할때)", map = 270050000 },
	{ text = "#r[핑크빈]#k #b신들의 황혼#k #e(팅길때)#n", map = 270050100 },
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

local ROUTES = {
	["0"] = MAIN,
	["1"] = TOWNS,
	["2"] = HUNTS,
	["3"] = BOSSES,
	["4"] = PARTIES,
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
	me:quest(875000):record("0")
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
		local route_q = me:quest(875000)
		local flag_q = me:quest(875001)
		if not route_q:started() then
			route_q:start("0")
		end
		if not flag_q:started() then
			flag_q:start("0")
		end
		if flag_q:record() == "0" then
			route_q:record("0")
		elseif flag_q:record() == "1" then
			flag_q:record("0")
		end

		local route = route_q:record()
		local entries = ROUTES[route]
		if entries == nil then
			return
		end
		local prompt = ""
		if route == "0" then
			prompt = "#k어느걸 이용할거야?"
		end
		choose(me, npc, prompt, entries)
	end
}

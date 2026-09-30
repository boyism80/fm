-- NPC name (String.wz/Npc.img.xml): 몽

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
	{ text = "#r[추천 Lv 130~150]#k 성벽 #b(그레이 벌쳐)#k", map = 211060410 },
	{ text = "#r[추천 Lv 150~160]#k 암벽 거인 #b(호넷)#k", map = 240091400 },
	{ text = "#r[추천 Lv 160~170]#k 레헬른 #b(닭고기)#k", map = 450003310 },
	{ text = "#r[추천 Lv 170~200]#k 아르카나 #b(부조화의 정령)#k", map = 450005431 },
	{ text = "#r[추천 Lv 200~250~]#k 미래의문 #b(돌연변이 슬라임)#k", map = 271010300 },
}

local BOSSES = {
	{ text = "#r[파풀라투스]#k #b시계탑 깊은 곳#k (입장)", map = 220080000 },
	{ text = "#r[크렉셀]#k #b사유지#k #e(점프맵)#n", map = 541020700 },
	{ text = "#r[자쿰]#k #b자쿰의 제단 입구#k (입장)", map = 211042300 },
	{ text = "#r[도로시]#k #b도로시의 영역#k #e(보스룸)#n", map = 123456788 },
	{ text = "#r[아인크라드]#k #b어딘가의 보스룸#k (소환할때)", map = 864000100 },
	{ text = "#r[마왕 발록]#k #b신전으로 가는길#k #k(입장)", map = 105100000 },
	{ text = "#r[텐구]#k #b월하죽림#k #e(나막신)#n", map = 800020130 },
	{ text = "#r[혼테일]#k #b생명의 동굴#k #k(입장)#n", map = 240050400 },
	{ text = "#r[보스웨이브]#k #b로비#k #k(입장)", map = 123456771 },
}

local PARTIES = {
	{ text = "#r[파티 3~6인]#k 레벨제한 : 10~250 #b월묘 파퀘#k", map = 100000200 },
	{ text = "#r[파티 4~6인]#k 레벨제한 : 21~250 #b커닝 파퀘#k", map = 103000000 },
	{ text = "#r[파티 2~6인]#k 레벨제한 : 30~250 #b몬스터 카니발#k", npc = 2042002 },
	{ text = "#r[파티 6~6인]#k 레벨제한 : 35~250 #b루디 파퀘#k", map = 221024500 },
	{ text = "#r[파티 6~6인]#k 레벨제한 : 51~250 #b올비 파퀘#k", map = 200080101 },
	{ text = "#r[파티 3~6인]#k 레벨제한 : 45~250 #b엘린숲 파퀘#k", map = 300030100 },
	{ text = "#r[파티 3~6인]#k 레벨제한 : 55~250 #b데비존 파퀘#k", map = 251010404 },
	{ text = "#r[파티 4~6인]#k 레벨제한 : 71~250 #b줄리엣 파퀘#k", map = 261000021 },
	{ text = "#r[파티 4~6인]#k 레벨제한 : 71~250 #b로미오 파퀘#k", map = 261000011 },
}

local MAIN = {
	{ text = "#k엠작 하기#k", npc = 9001130 },
	{ text = "#k마나 하락#n#k", npc = 3002105 },
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

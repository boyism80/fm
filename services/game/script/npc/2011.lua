-- NPC name (String.wz/Npc.img.xml): 할로 캣

local TOWNS = {
	{ text = "자유시장", map = 123456789 },
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
	{ text = "#r[혼테일]#k #b생명의 동굴#k #k(입장)", map = 240050400 },
}

local PARTIES = {
	{ text = "#r홍보 상점#k", npc = 9075005 },
	{ text = "#b일지 상점#k", npc = 9000560 },
	{ text = "#b마북 상점#k", npc = 9000032 },
}

local MAIN = {
	{ text = "티어 이용#r//승급,퀘스트 받기#k", npc = 2009 },
	{ text = "위치 이동#r//사냥터,보스 이동#k", npc = 2010 },
	{ text = "상점 이용#r//부산물 아이템 교환#k", npc = 2012 },
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
			prompt = "#k반복 퀘스트 진행이 어려우시다구요? #e@이동#n 을 입력후 #e강화소 이동#n 을 이용하여 티어 엔피시 에게 받아주시면 더욱 편한 시스템 이용이 가능합니다."
		end
		choose(me, npc, prompt, entries)
	end
}

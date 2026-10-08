-- NPC name (String.wz/Npc.img.xml): 차원의 거울

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

local MAIN = {
	{ text = "#k마을 이동#k", menu = TOWNS },
	{ text = "#k거울 이동#n#k", npc = 9001013 },
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

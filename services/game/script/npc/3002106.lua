-- NPC name (String.wz/Npc.img.xml): 랑

local MAIN = {
	{ text = "#r#e누적 보상#k", npc = 2017 },
	{ text = "#r#e강화 이전#k", npc = 9000193 },
	{ text = "#r#e붕어 가챠#k", npc = 9000601 },
	{ text = "#r#e핫탐 가챠#k#n", npc = 2016 },
	{ text = "#k캐시 정리#k", npc = 9000556 },
	{ text = "#k최신 캐시#k", npc = 9000562 },
	{ text = "#k성형 하기#k", npc = 9000555 },
	{ text = "#k출석 체크#k", npc = 9001108 },
	{ text = "#k시그 정축#k", npc = 9001137 },
	{ text = "#k모험 정축#k", npc = 1530731 },
	{ text = "#k만렙 스킬#k", npc = 3002102 },
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
		choose(me, npc, "#k너에게 #e출석 체크#n와 #e누적 보상 시스템#n을 추천할게!", MAIN)
	end
}

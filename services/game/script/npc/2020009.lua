-- NPC name (String.wz/Npc.img.xml): 로베이라

local APPROVAL_QUEST = 100000
local MIN_LEVEL = 50

local function grant_approval(me, npc)
	local q = me:quest(APPROVAL_QUEST)
	if me:level() >= MIN_LEVEL and q ~= nil and not q:started() and not q:completed() then
		me:dialog(npc, "자쿰던전 퀘스트를 허가해 달라 그거군... #b아도비스#k인가... 아무튼 좋네! 자네라면 그 던전을 탐색하는 데 모자람이 없겠지. 그럼 아무쪼록 조심하길 바라네.")
		if q:wz() == nil then
			q:start("")
		else
			q:start(0, true)
		end
		return
	end
	if q ~= nil and (q:started() or q:completed()) then
		me:dialog(npc, "자네는 이미 자쿰던전 퀘스트를 허가받지 않았는가? 아무쪼록 조심하길 바라네.")
		return
	end
	me:dialog(npc, "아직 자네의 실력이 부족해 보이는군. 레벨 50 이상이 된 후에 다시 찾아오게나.")
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "나에게 무슨 볼일이라도 있는가?", {
			"3차 전직을 하고 싶습니다.",
			"자쿰던전 퀘스트를 허가해 주세요",
		})
		if sel == 1 then
			me:dialog(npc, "3차 전직은 아직 준비 중이네. 조금 더 기다려 주게.")
		elseif sel == 2 then
			grant_approval(me, npc)
		end
	end
}

-- NPC name (String.wz/Npc.img.xml): 델리

local QUEST = 6410
local PROGRESS_QUEST = 6411
local ROAD_MAP = 925010200
local PROTECT_MAP = 925010300
local EXIT_MAP = 120000104
local PROTECT_MS = 360000

local function protect(me, npc)
	if me:dialog(npc, "당신은 누구시죠..? 네..? 슈린츠가 보내서 왔다구요?", false, true) == false then
		return
	end
	if me:dialog_yes_no(npc, "다행이다.. 정말 무서웠어요. 그런데 아직도 저를 노리는 무서운 몬스터들이 많이 있어요. 그 몬스터들을 물리쳐 주시면 좋겠어요.") == false then
		me:dialog(npc, "... 실망이에요.")
		return
	end
	if me:dialog(npc, "정말 감사해요, 6분 동안 몬스터들로부터 저를 보호해 주시면 된답니다.", false, true) == false then
		return
	end
	local sm = me:state_machine()
	if sm == nil then
		return
	end
	sm:set_property("protect", "1")
	sm:restart_timer(PROTECT_MS)
	me:map(sm:map(PROTECT_MAP), 0)
end

local function thank(me, npc)
	if me:quest(QUEST):started() == false then
		me:map(EXIT_MAP)
		return
	end
	if me:dialog(npc, "저를 구해주셔서 정말 고마워요. 당신은 친절한 사람이군요.", false, true) == false then
		return
	end
	me:quest(PROGRESS_QUEST):start(npc, "p2")
	me:show_quest_completion(QUEST)
end

return {
	on_click = function(me, npc)
		if me:map():template_id() == ROAD_MAP then
			protect(me, npc)
			return
		end
		thank(me, npc)
	end
}

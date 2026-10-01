-- NPC name (String.wz/Npc.img.xml): 조교 챙

local QUEST = 3239
local PART = 4031092
local PART_COUNT = 10
local FACTORY_MAP = 220020000
local FACTORY_PORTAL = 2
local REWARDS = { 2040704, 2040705, 2040707, 2040708 }

local function part_count(me)
	local count = 0
	for _, it in pairs(me:item(PART)) do
		count = count + it:count()
	end
	return count
end

local function enter(me, npc)
	local quest = me:quest(QUEST)
	if quest:completed() then
		me:dialog(npc, "아, 자넨 저번에 정비부품 찾는걸 도와줬던 친구아닌가? 그땐 정말 고마웠다네.")
		return
	end
	if quest:started() == false then
		me:dialog(npc, "흐음.. 이곳엔 아무나 들어갈 수 없다네.")
		return
	end
	if me:dialog(npc, "오, 정비 부품 찾는걸 도와주러 왔는가? 이곳으로 들어가보면 상자가 있을 걸세. 상자들을 부수면 찾을 수 있을거야.", false, true) == false then
		return
	end
	me:rmitem(PART)
	local sm, err = state_machine("machine_room"):start_solo(me)
	if sm == nil then
		log("machine_room start_solo:", err)
		me:dialog(npc, "이미 이 안에 다른 플레이어가 들어간 것 같군. 잠시 후에 다시 시도해 보게나.")
	end
end

local function turn_in(me, npc)
	local count = part_count(me)
	if count < PART_COUNT then
		if me:dialog_yes_no(npc, "흐음. 아직 기계 부품 10개를 모으지 못한 것 같구만. 기계 부품 모으기는 중단하고 여기서 나가보겠는가?") == false then
			me:dialog(npc, "좋네. 조금 더 노력해 주면 고맙겠네.")
			return
		end
		me:map(FACTORY_MAP, FACTORY_PORTAL)
		return
	end
	if me:dialog(npc, "오! 기계 부품 10개를 모두 모아왔군! 정말 고맙네. 여기 내 보답을 받게. 받기 전에 인벤토리 공간이 충분한지 확인해주게나~", false, true) == false then
		return
	end
	local reward = REWARDS[math.random(#REWARDS)]
	if me:exchange({ item = { [PART] = count } }, { item = { [reward] = 1 }, exp = 2700 }) ~= ExchangeResult.OK then
		me:dialog(npc, "자네.. 소비 인벤토리 공간이 부족한 것 같은데? 다시 확인해보게나.")
		return
	end
	me:quest(QUEST):force_complete(npc)
	me:map(FACTORY_MAP, FACTORY_PORTAL)
end

return {
	on_click = function(me, npc)
		if me:map():template_id() == FACTORY_MAP then
			enter(me, npc)
			return
		end
		turn_in(me, npc)
	end
}

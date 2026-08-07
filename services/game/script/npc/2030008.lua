-- NPC name (String.wz/Npc.img.xml): 아도비스

local pq = require("script/lib/party_quest")

local GROUP_NAME = "zakum_party_quest"
local MIN_LEVEL = 50
local STAGE2_MAP = 280020000
local APPROVAL_QUEST = 100000
local STAGE1_QUEST = 100001
local STAGE2_QUEST = 100002
local STAGE3_QUEST = 100003
local FIRE_ORE_PIECE = 4031061
local VOLCANO_BREATH = 4031062
local HECTOR_TAIL = 4000051
local EYE_OF_FIRE = 4001017
local REINFORCED_BOTTLE = 4001109
local FIRE_DEMON = "fire_demon"

local function ensure_quest_started(me, quest_id, record)
	local q = me:quest(quest_id)
	if q == nil then
		return nil
	end
	if q:started() or q:completed() then
		return q
	end
	if record == nil then
		record = ""
	end
	if q:wz() == nil then
		q:start(record)
	else
		q:start(0, true)
	end
	return me:quest(quest_id)
end

local function quest_record_starts(me, quest_id, prefix)
	local q = me:quest(quest_id)
	if q == nil then
		return false
	end
	local rec = q:record()
	if rec == nil or rec == "" then
		return false
	end
	return string.sub(rec, 1, #prefix) == prefix
end

local function set_quest_record(me, quest_id, value)
	local q = me:quest(quest_id)
	if q == nil then
		return false
	end
	if not q:started() and not q:completed() then
		ensure_quest_started(me, quest_id, value)
		q = me:quest(quest_id)
	end
	if q == nil then
		return false
	end
	return q:record(value)
end

local function start_fire_demon(me, npc)
	local group = state_machine(FIRE_DEMON)
	if group == nil then
		me:dialog(npc, "오류가 발생했습니다.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이미 다른 사람이 용암의 심장부에 들어가 있습니다.")
		return
	end
	local sm, err = group:start_solo(me)
	if sm == nil then
		me:dialog(npc, "오류가 발생했습니다.")
		if err ~= nil then
			log("fire_demon start_solo:", err)
		end
	end
end

local function start_zakum_pq(me, npc)
	local party = me:party()
	if party == nil or party:leader_id() ~= me:id() then
		me:dialog(npc, "퀘스트를 받아온 후 파티를 맺어 파티장이 나에게 말을 걸면 여러가지 퀘스트를 수행할 수 있다네. 모든 준비가 끝났다면 파티장한테 나에게 말을 걸라고 전해주게나.")
		return
	end
	if not pq.party_all_here(me) then
		me:dialog(npc, "모든 파티원이 다 여기에 있어야 퀘스트를 수행할 수가 있다네. 모든 준비가 끝났다면 파티장한테 나에게 말을 걸라고 전해주게나.")
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local ok = true
	for _, mem in ipairs(party:members()) do
		if mem == nil or mem:level() < MIN_LEVEL then
			ok = false
			break
		end
	end
	if not ok then
		me:dialog(npc, "모든 파티원이 다 여기에 있어서 퀘스트를 진행 할 수가 있다네.")
		return
	end
	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "파티퀘스트에 현재 문제가 생긴것 같군. 미안하지만 지금은 입장할 수 없겠네.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이미 다른 파티가 안으로 들어가 퀘스트 클리어에 도전하고 있습니다.")
		return
	end
	ensure_quest_started(me, STAGE1_QUEST, "")
	local sm, err = group:start_party(me, party)
	if sm == nil then
		me:dialog(npc, "지금은 무언가 문제가 생긴 것 같군..")
		if err ~= nil then
			log("zakum_party_quest start_party:", err)
		end
	end
end

local function handle_stage2(me, npc)
	local stage1 = me:quest(STAGE1_QUEST)
	local stage2 = me:quest(STAGE2_QUEST)
	local stage1_done = stage1 ~= nil and stage1:completed()
	local stage2_done = stage2 ~= nil and stage2:completed()
	local cleared = quest_record_starts(me, APPROVAL_QUEST, "Zakum1Clear")
	if not stage1_done and not stage2_done and not cleared then
		me:dialog(npc, "자네는 1단계를 진행중인 것 같군 그래. 2단계를 도전하기 위해서는 1단계를 성공적으로 클리어 한 상태여야만 하네. 우선 1단계를 클리어 하게나.")
		return
	end
	local msg
	if stage2_done then
		msg = "흠... 자네는 이미 2단계를 클리어 한 적이 있는 모양이로군. 하지만 원한다면 언제든지 또 도전이 가능하다네. 어떤가... 다시 한 번 2단계에 도전해 보겠는가?"
	else
		msg = "자네는 1단계를 무사히 클리어 했군 그래. 하지만 아직 자쿰 던전 보스를 만나기 위해서는 많은 난관이 남아있지. 어떤가... 2단계에 도전해 보겠는가?"
	end
	if not me:dialog_yes_no(npc, msg) then
		me:dialog(npc, "그렇군... 하지만 언제라도 마음이 정해졌다면 다시 날 찾아와 주길 바라네.")
		return
	end
	if not me:dialog(npc, "좋네! 이제부터 자네를 수 많은 장애물들이 있는 맵으로 이동될 것일세. 그곳의 가장 안쪽에는 보물상자가 있는데 보물상자를 조사하면 보스를 소환하는 데 필요한 아이템의 재료 중 하나를 얻을 수 있을 거야. 재료를 얻어서 나에게 가져와 주게나. 그럼 힘내주게!", true, true) then
		return
	end
	if stage2 == nil or (not stage2:started() and not stage2:completed()) then
		ensure_quest_started(me, STAGE2_QUEST, "")
	end
	me:map(STAGE2_MAP)
end

local function handle_refine(me, npc)
	local stage2 = me:quest(STAGE2_QUEST)
	if stage2 == nil or not stage2:completed() then
		me:dialog(npc, "아직 자네는 이전 단계를 클리어 하지 않은 것 같군. 이전 단계를 클리어 한 후에 다시 시도해 주길 바라네.")
		return
	end
	local stage3 = me:quest(STAGE3_QUEST)
	local stage3_done = stage3 ~= nil and stage3:completed()
	local cleared2 = quest_record_starts(me, APPROVAL_QUEST, "Zakum2Clear")
	local msg
	if stage3_done then
		msg = "흠... 자네는 일전에 #b불의 눈#k을 제련해 간 사람이 아닌가. 그런데 나에게 무슨 볼일인가. 다시 한 번 #b불의 원석 조각#k과 #b화산의 숨결#k을 조합하여 #b불의 눈#k을 만들어 볼텐가?"
	else
		msg = "자네는 2단계를 무사히 클리어 했군 그래. 하지만 아직 자쿰 던전 보스를 만나기 위해서는 마지막 과정이 남아 있다네. 어떤가... 지금까지 모아온 재물을 조합해서 제련을 해 볼 텐가?"
	end
	if not me:dialog_yes_no(npc, msg) then
		me:dialog(npc, "그렇군... 하지만 언제라도 마음이 정해졌다면 다시 날 찾아와 주길 바라네.")
		return
	end
	if not cleared2 then
		me:dialog(npc, "흠, #b불의 원석 조각#k와 #b화산의 숨결#k을 조합하면 보스를 부르는데 제물로 바쳐야 하는 아이템인 #b불의 눈#k를 만들 수 있지. 하지만... 쿨럭쿨럭...! 보다시피 몸이 좋질 않아 뭐든 구하기가 여간 어려운게 아니라네. 그러니 혹, #b헥터의 꼬리 30개#k만 구해다 줄 소 있소? 어디에 쓸 지는 묻지 말고... 흠흠...")
		set_quest_record(me, APPROVAL_QUEST, "Zakum2Clear")
		return
	end
	if not (pq.has_item(me, FIRE_ORE_PIECE)
		and pq.has_item(me, VOLCANO_BREATH)
		and pq.has_item(me, HECTOR_TAIL, 30)) then
		me:dialog(npc, "아직 #b헥터의 꼬리 30개#k를 구해오지 못한 것 같군 그래. 이것만 모아와 준다면 그동안 자네들이 가져다 준 것을 제련해서 특별한 물건을 만들수도 있을 것 같은데 말이지. 아참, 이전 단계에서 모은 #b불의 원석#k과 #b화산의 숨결#k도 물론 모아와야 한다는 것을 잊지 말게.")
		return
	end
	if not me:dialog(npc, "하하하, 전광석화와 같이 금세 만들어주겠네.", true, true) then
		return
	end
	local code = me:exchange({
		item = {
			[FIRE_ORE_PIECE] = 1,
			[VOLCANO_BREATH] = 1,
			[HECTOR_TAIL] = 30,
		},
	}, {
		item = { [EYE_OF_FIRE] = 5 },
	})
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "인벤토리 공간을 확인한 뒤 다시 말을 걸어 주게.")
		return
	end
	ensure_quest_started(me, STAGE3_QUEST, "")
	local q3 = me:quest(STAGE3_QUEST)
	if q3 ~= nil and q3:started() then
		q3:force_complete(npc)
	end
	set_quest_record(me, APPROVAL_QUEST, "")
	me:dialog(npc, "여깄네. 자, 이제 저 문이 열리면 안으로 들어갈 수 있을 걸세. 이 #b불의 눈#k을 갖고 있어야만 문을 통해 입장을 할 수 있다네. 뭐, 열 명까지 입장이 가능하다던가...")
end

local function handle_explain(me, npc)
	if not me:dialog(npc, "어떻게 해야하는지 감이 잡히질 않는 모양이로구만? 이 퀘스트를 하려면 장로의 관저에 가서 퀘스트 수행 승인을 받아와야 하네. 그렇지 않고서 들여보냈다가 추궁 받긴 싫다구... 퀘스트를 받아온 이들로만 파티를 맺으면 들여보낼 수 있네.", true, true) then
		return
	end
	if not me:dialog(npc, "순서대로 수행해나가면 자쿰던전의 보스를 만날 수 있을 걸세. 내가 모아오라는 것을 모아오면 제물을 만들어주겠네. 그 제물을 제단에 바치라고. 그럼 당신이 원하는 걸 보게 될 것이네. 그러기 위해서는 먼저 폐광 동굴을 조사하여 #b불의 원석#k을 가져와야 한다네.", true, true) then
		return
	end
	if not me:dialog(npc, "거기는 #b불의 원석#k 외에도 종이문서를 발견할 수 있을 걸세. 그걸 #b아우라#k에게 가져다주면 도움이 되는 것을 만들어줄 거야. 다음은 용암지대를 건너가 #b화산의 숨결#k을 구해와야 한다네. 험난한 길이겠지만... 제물을 만들기 위해서는 반드시 필요한 것이라네.", true, true) then
		return
	end
	me:dialog(npc, "#b화산의 숨결#k까지 구했다면 1, 2단계에서 얻은 #b불의 원석 조각#k과 #b화산의 숨결#k을 제련하는 일이 필요하다네. 제련은 내가 할 수 있으니 걱정 말게. 이 모든게 끝났다면 자쿰의 보스를 만나는 일만 남게 되지. 결코 쉽지는 않겠지만 열심히 해 보게나.")
end

return {
	on_click = function(me, npc)
		local approved = false
		local q = me:quest(APPROVAL_QUEST)
		if q ~= nil and (q:started() or q:completed()) then
			approved = true
		end
		local has_bottle = pq.has_item(me, REINFORCED_BOTTLE)
		local options = {}
		local actions = {}
		if approved then
			options[#options + 1] = "폐광 동굴을 조사하러 떠난다. (1단계)"
			actions[#actions + 1] = "pq"
			options[#options + 1] = "자쿰 던전을 탐사한다. (2단계)"
			actions[#actions + 1] = "stage2"
			options[#options + 1] = "제련을 요청한다. (3단계)"
			actions[#actions + 1] = "refine"
			options[#options + 1] = "퀘스트에 관한 설명을 듣는다."
			actions[#actions + 1] = "explain"
		end
		if has_bottle then
			options[#options + 1] = "용암의 심장부로 들어간다."
			actions[#actions + 1] = "fire"
		end
		if #options == 0 then
			me:dialog(npc, "이 앞은 혼자서는 결코 해결할 수 없는 미궁으로 가득차 있지. 혹시라도 도전해 보고 싶다면 엘나스의 장로의 관저에 있는 각 직업별 장로에게 퀘스트를 받아오도록 하게.")
			return
		end
		local sel = me:dialog_list(npc, "뭐... 좋소. 당신들은 충분한 자격이 되어 보이는군. 어느 단계에 도전해 보겠소?", options)
		if sel == nil or sel < 1 or sel > #actions then
			me:dialog(npc, "그렇군... 하지만 언제라도 마음이 정해졌다면 다시 날 찾아와 주길 바라네.")
			return
		end
		local action = actions[sel]
		if action == "pq" then
			start_zakum_pq(me, npc)
		elseif action == "stage2" then
			handle_stage2(me, npc)
		elseif action == "refine" then
			handle_refine(me, npc)
		elseif action == "explain" then
			handle_explain(me, npc)
		elseif action == "fire" then
			start_fire_demon(me, npc)
		end
	end
}

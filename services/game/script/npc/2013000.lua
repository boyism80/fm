-- NPC name (String.wz/Npc.img.xml): 요정 웡키

local pq = require("script/lib/party_quest")

local GROUP_NAME = "orbis_party_quest"
local MIN_PARTY_SIZE = 4
local MIN_LEVEL = 51
local SCALE_LEVEL = 70
local FEATHER = 4001158
local BRACELET = 1082232
local ENTRY_MAP = 200080101
local GIVEUP_MAP = 920010000

local function try_start(me, npc)
	local party = me:party()
	if party == nil then
		me:dialog(npc, "파티원을 모으거나 파티에 참가한 후에 여신의 흔적을 찾으러 갈 수 있어.")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "흐음. 너는 파티장이 아닌 것 같은데? 퀘스트를 시작하려면 파티장이 내게 말을 걸어야 한다구.")
		return
	end
	local map = me:map()
	if map == nil or map:wz() == nil then
		return
	end
	local map_id = map:wz().id
	local ok = true
	local in_map = 0
	local has_gm = false
	for _, mem in ipairs(party:members()) do
		if mem ~= nil then
			if mem:level() < MIN_LEVEL then
				ok = false
			end
			if mem:map_id() == map_id then
				local ch = map:characters()[mem:id()]
				if ch ~= nil then
					if pq.is_gm(ch) then
						has_gm = true
					end
					in_map = in_map + 1
				else
					ok = false
				end
			end
		end
	end
	if not ok or (in_map ~= MIN_PARTY_SIZE and not has_gm) then
		me:dialog(npc, "파티 요구 조건이 맞지 않는 것 같은데? 파티는 #r레벨 51 이상의 파티원으로만 이루어진 4명의 파티#k만 입장이 가능해. 다시 한번 확인해 봐.")
		return
	end
	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이미 다른 파티가 이 안에서 퀘스트에 도전하는 중이야. 잠시 후 다시 시도해 주거나 채널을 변경해 봐.")
		return
	end
	local sm, err = group:start_party(me, party, SCALE_LEVEL)
	if sm == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		if err ~= nil then
			log("orbis_party_quest start_party:", err)
		end
	end
end

local function claim_bracelet(me, npc)
	if not pq.has_item(me, FEATHER, 10) then
		me:dialog(npc, "#b여신의 깃털 10개#k는 가지고 있는거야? 여신의 아이템을 받고 싶으면 #b여신의 깃털 10개#k를 가지구 오라구~ 또 만들어 줄게!")
		return
	end
	local code = me:exchange(
		{ item = { [FEATHER] = 10 } },
		{ item = { [BRACELET] = 1 } }
	)
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간을 확보한 뒤 다시 말을 걸어줘.")
		return
	end
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "#b여신의 깃털 10개#k는 가지고 있는거야? 여신의 아이템을 받고 싶으면 #b여신의 깃털 10개#k를 가지구 오라구~ 또 만들어 줄게!")
		return
	end
	me:dialog(npc, "고마워~~! 여기 #b여신의 팔찌#k를 받아~ 혹시라도 또 필요하면 #b여신의 깃털 10개#k를 가지구 오라구~ 또 만들어 줄게!")
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		if map_id == GIVEUP_MAP then
			if not pq.is_leader(me) then
				me:dialog(npc, "아직 포기하긴 이르다구? 조금 더 열심히 노력해 봐~")
				return
			end
			if me:dialog_yes_no(npc, "음.. 정말 이곳에서 나가고 싶어? 여기서 나간다면 처음부터 다시 시작해야만 해.") then
				me:map(920011200)
			else
				me:dialog(npc, "아직 포기하긴 이르다구? 조금 더 열심히 노력해 봐~")
			end
			return
		end
		if map_id ~= ENTRY_MAP then
			return
		end
		local sel = me:dialog_list(npc, "#e<파티퀘스트 : 여신의 흔적>#n\r\n이야아~ 반가워! 여신의 탑을 모험하고 싶어?\r\n\r\n#b", {
			"입장을 신청한다.",
			"여신의 탑에 대해 묻는다.",
			"여신의 깃털을 다른 아이템과 바꾼다.",
		})
		if sel == 1 then
			try_start(me, npc)
		elseif sel == 2 then
			me:dialog(npc, "며칠 전에 갑자기 오르비스 구름 위에 가장 아름다운 여신, 미네르바와 신비한 탑이 나타났어.")
			me:dialog(npc, "그런데 어느 순간 갑자기 여신은 사라지고 탑은 제어 불능에 빠져버렸어. 아무래도 여신이 탑 어딘가에 갇혀 있는 것 같아. 너는 어떻게 생각해? 그 탑을 탐험하고 여신의 흔적을 찾아 떠나보고 싶지 않아?")
		elseif sel == 3 then
			local craft = me:dialog_list(npc, "나의 배고픔을 달래 주던 #b#h ##k님이잖아! 그동안 고마웠다는 표시로 #b여신의 깃털을 여신의 아이템으로#k 만들어 줄게. 뭘 만들어 줄까?#b\r\n", {
				"여신의 팔찌 #r(여신의 깃털 10개 필요)",
			})
			if craft == 1 then
				claim_bracelet(me, npc)
			end
		end
	end
}

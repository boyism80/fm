-- State machine (old/src/server/marriage/MarriageEventAgent.java): 결혼식

local LOBBY = 680000200
local CATHEDRAL = 680000210
local PHOTO = 680000300
local PARTY = 680000400
local BONUS = 680000401
local EXIT_MAP = 680000500
local WAIT_MS = 600000
local CEREMONY_MS = 600000
local PHOTO_MS = 300000
local PARTY_MS = 180000
local BONUS_MS = 180000
local VOW_EFFECT = 5120025
local CHEAP_TICKET = 5251004
local SWEETIE_TICKET = 5251005

local BUFFS = {
	[5251004] = 2022196,
	[5251005] = 2022197,
	[5251006] = 2022200,
}

local VOWS = {
	{ at = 0, text = "오늘 우리는 두 젊은이를 축복하기 위해 모였습니다." },
	{ at = 10000, text = "흔히 하늘이 내린 인연은 리본돼지의 붉은 리본으로 이어져 있다고 합니다." },
	{ at = 20000, text = "제가 보기에는 이 두사람은 그 리본의 끝에 있는 서로를 찾아낸 것 같습니다." },
	{ at = 30000, text = "수 많은 여행자들 중에서 서로를 찾아낸 두 사람이야말로 진정한 행운아일겁니다." },
	{ at = 40000, text = "서로에게 찾아 온 이 행운을 행복으로 만들어가는 것이 앞으로 두 사람에게 주어진 퀘스트입니다." },
	{ at = 50000, text = "신랑, 검은머리가 예티의 털처럼 하얗게 변할때까지 신부를 사랑하시겠습니까?" },
	{ at = 60000, text = "신부, 엘나스산의 만년설이 모두 녹아 니할사막처럼 될 때까지 신랑을 사랑하시겠습니까?" },
	{ at = 80000, text = "오늘 모인 하객분들이 두 사람의 맹세에 대한 증인이 되어주십시오." },
	{ at = 90000, text = "모두의 축복 속에 두 젊은이가 부부가 되었음을 선포합니다." },
	{ at = 100000, text = "신랑은 신부에게 키스해도 좋습니다." },
}
local CHAPEL_AT = 110000

local function couple(sm, player)
	local id = tostring(player:id())
	return id == sm:get_property("groom_id") or id == sm:get_property("bride_id")
end

local function finale(sm)
	local level = tonumber(sm:get_property("finale"))
	local ticket = tonumber(sm:get_property("ticket"))
	if level == 0 then
		if ticket == CHEAP_TICKET then
			sm:finish(EXIT_MAP)
			return
		end
		sm:cancel("vow")
		sm:warp_all(CATHEDRAL, PHOTO)
		sm:set_property("finale", "1")
		sm:start_timer(PHOTO_MS)
	elseif level == 1 then
		sm:warp_all(PHOTO, PARTY)
		if ticket == SWEETIE_TICKET then
			sm:set_property("finale", "3")
		else
			sm:set_property("finale", "2")
		end
		sm:start_timer(PARTY_MS)
	elseif level == 2 then
		for _, player in pairs(sm:map(PARTY):characters()) do
			if couple(sm, player) then
				player:map(sm:map(BONUS))
			else
				player:map(EXIT_MAP)
			end
		end
		sm:set_property("finale", "3")
		sm:start_timer(BONUS_MS)
	else
		sm:finish(EXIT_MAP)
	end
end

return {
	on_init = function(group)
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		sm:set_property("state", "waiting")
		sm:set_property("finale", "0")
		return { LOBBY, CATHEDRAL, PHOTO, PARTY, BONUS }
	end,

	on_start = function(sm)
		sm:start_timer(WAIT_MS)
		sm:message("제한시간 내에 웨딩홀 내 안나 수녀님을 통해 결혼식을 시작해야 이벤트와 각종 보상을 얻을 수 있습니다.", Msg.PinkText)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(LOBBY))
	end,

	on_player_leave = function(sm, player, reason)
		if couple(sm, player) then
			sm:finish(EXIT_MAP)
		end
	end,

	on_enter_cathedral = function(sm)
		if sm:get_property("state") ~= "waiting" then
			return
		end
		for _, player in pairs(sm:map(LOBBY):characters()) do
			if couple(sm, player) then
				player:map(sm:map(CATHEDRAL), 2)
			else
				player:map(sm:map(CATHEDRAL), 1)
			end
		end
	end,

	on_ceremony = function(sm)
		if sm:get_property("state") ~= "waiting" then
			return
		end
		sm:set_property("state", "ceremony")
		sm:set_property("vow", "0")
		sm:map(CATHEDRAL):music("BgmGL/cathedral")
		sm:start_timer(CEREMONY_MS)
		local buff = BUFFS[tonumber(sm:get_property("ticket"))]
		for _, player in pairs(sm:map(CATHEDRAL):characters()) do
			player:use_item(buff)
		end
		sm:call_hook("on_vow")
	end,

	on_vow = function(sm)
		local hall = sm:map(CATHEDRAL)
		local index = tonumber(sm:get_property("vow")) + 1
		local vow = VOWS[index]
		if vow == nil then
			hall:music("BgmGL/chapel")
			hall:weather(0)
			return
		end
		sm:set_property("vow", tostring(index))
		hall:weather(VOW_EFFECT, vow.text)
		hall:yellow_chat(vow.text)
		if index == #VOWS then
			hall:wedding_effect(tonumber(sm:get_property("groom_id")), tonumber(sm:get_property("bride_id")))
			if finish_wedding(tonumber(sm:get_property("marriage_id"))) == false then
				log("wedding finish_wedding failed:", sm:get_property("marriage_id"))
			end
			sm:set_property("state", "kissed")
			sm:after("vow", CHAPEL_AT - vow.at, "on_vow")
			return
		end
		sm:after("vow", VOWS[index + 1].at - vow.at, "on_vow")
	end,

	on_finale = function(sm)
		if sm:get_property("state") ~= "kissed" then
			return
		end
		finale(sm)
	end,

	on_scheduled_timeout = function(sm)
		if sm:get_property("state") == "kissed" then
			finale(sm)
			return
		end
		sm:message("제한시간 내에 결혼식을 시작하지 않아 결혼식이 취소되었습니다.", Msg.PinkText)
		sm:finish(EXIT_MAP)
	end,
}

-- NPC name (String.wz/Npc.img.xml): 안젤리크

local LOBBY = 680000200

local function wishes(text)
	local out = {}
	for wish in string.gmatch(text, "[^\n]+") do
		table.insert(out, wish)
	end
	return out
end

local function give(me, npc)
	local sm = me:state_machine()
	if sm == nil then
		return
	end
	local id = tostring(me:id())
	if id == sm:get_property("groom_id") or id == sm:get_property("bride_id") then
		me:dialog(npc, "결혼을 진심으로 축하드려요~ 하객 분들은 저를 통해 결혼 선물을 여러분들께 드릴 수 있답니다. 결혼을 마치고 나가는 길에 저를 통해 선물을 수령하실 수 있으니 걱정 마세요~")
		return
	end
	local selected = me:dialog_list(npc, "어서오세요. 신랑 신부에게 주고 싶은 선물은 제가 대신 관리하고 있어요. 어느쪽 하객분이신가요?", {
		"신랑에게 선물을 주고 싶습니다.",
		"신부에게 선물을 주고 싶습니다.",
	})
	if selected == 1 then
		me:open_wedding_gift(tonumber(sm:get_property("groom_id")), wishes(sm:get_property("groom_wishes")))
	elseif selected == 2 then
		me:open_wedding_gift(tonumber(sm:get_property("bride_id")), wishes(sm:get_property("bride_wishes")))
	end
end

local function receive(me, npc)
	local marriage = me:marriage()
	if marriage == nil then
		me:dialog(npc, "저는 결혼 선물을 대신 관리하고 있답니다.")
		return
	end
	if not marriage:married() then
		me:dialog(npc, "선물을 받으실 수 있는 상태가 아닌 것 같군요.")
		return
	end
	if not me:open_wedding_gift_box() then
		me:dialog(npc, "받으실 수 있는 선물이 없군요.")
	end
end

return {
	on_click = function(me, npc)
		if me:map():wz():id() == LOBBY then
			give(me, npc)
			return
		end
		receive(me, npc)
	end
}

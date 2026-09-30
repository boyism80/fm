local JAIL_MAP = 180000002
local MAX_TEXT_BYTES = 60

local M = {}

-- The client limits megaphone text in EUC-KR bytes: ASCII is 1 byte, Korean is 2.
local function length(text)
	local _, ascii = string.gsub(text, "[\1-\127]", "")
	local _, multi = string.gsub(text, "[\192-\255]", "")
	return ascii + multi * 2
end

local function medal(me)
	local item = me:equipped(EquipmentPart.Medal)
	if item == nil then
		return ""
	end
	local name = item2name(item:wz():id())
	if name == nil then
		return ""
	end
	return "<" .. string.gsub(name, "의 훈장", "") .. "> "
end

function M.send(me, text, msg_type, scope, ear)
	if me:map():wz():id() == JAIL_MAP then
		me:message("이곳에서 사용할 수 없습니다.", Msg.PinkText)
		return false
	end
	if get_megaphone_muted() then
		me:message("현재 확성기 사용 금지 상태입니다.", Msg.PinkText)
		return false
	end
	if text == nil or length(text) == 0 or length(text) > MAX_TEXT_BYTES then
		return false
	end
	return me:message(medal(me) .. me:name() .. " : " .. text, msg_type, scope, ear)
end

return M

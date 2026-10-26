local MAX_TEXT_BYTES = 60

local M = {}

-- A weather item shows its effect with the user's words to the map, and some also give everyone on it a buff.
function M.use(me, item_id, text)
    if text == nil or not valid_text(text, MAX_TEXT_BYTES) then
        return false
    end
    local map = me:map()
    if not map:start_weather(item_id, text) then
        me:message("이미 다른 날씨 효과가 진행 중입니다.", Msg.PinkText)
        return false
    end
    map:give_state_change(item_id)
    return true
end

return M

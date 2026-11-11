-- Item name (String.wz/Cash.img.xml): 홀리 심볼

local weather = require("script/lib/weather")

local HOLY_SYMBOL = 9001002
local DURATION_SECONDS = 3600

return {
    -- Besides the weather, everyone on the map gets an hour of Holy Symbol, as the original did.
    on_cash = function(me, item_id, text)
        if not weather.use(me, item_id, text) then
            return false
        end
        for _, ch in pairs(me:map():characters()) do
            local skill = ch:skill_at(HOLY_SYMBOL, 1)
            if skill ~= nil then
                ch:buff(skill, BuffFlag.HolySymbol, skill:effect().x, { time = DURATION_SECONDS })
            end
        end
        return true
    end,
}

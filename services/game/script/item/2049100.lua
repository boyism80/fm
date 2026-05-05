-- Item name (String.wz/Consume.img.xml): 혼돈의 주문서 60%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, nil, nil, function(target, scroll)
        local scroll_wz = scroll:wz()
        if scroll_wz == nil then
            return
        end
        local z = scroll_wz:rand_stat()
        if z == nil or z <= 0 then
            return
        end

        local keys = {
            "str", "dex", "int", "luk",
            "pad", "pdd", "mad", "mdd",
            "acc", "avoid", "speed", "jump",
            "max_hp", "max_mp"
        }
        local bonus = {}
        for _, key in ipairs(keys) do
            local total = target:total_stat(key)
            if total ~= nil and total > 0 then
                local delta = math.random(-z, z)
                if delta ~= 0 then
                    bonus[key] = delta
                end
            end
        end
        target:add_bonus_stats(bonus)
    end)
end

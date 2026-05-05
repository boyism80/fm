-- Item name (String.wz/Consume.img.xml): 전신 갑옷 방어력 주문서 15%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Top, nil, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end

-- Item name (String.wz/Consume.img.xml): 스태프 마력 주문서 30%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Staff, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end

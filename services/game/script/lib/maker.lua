local M = {}

local function recipe_text(recipe)
    if recipe.text ~= nil then
        return recipe.text
    end
    return "#z" .. recipe.item .. "##k" .. recipe.label
end

local function choose_recipe(me, npc, category)
    if category.confirm ~= nil then
        if not me:dialog_yes_no(npc, category.confirm) then
            return nil
        end
        return category.recipes[1]
    end

    local options = {}
    for i, recipe in ipairs(category.recipes) do
        options[i] = recipe_text(recipe)
    end
    local sel = me:dialog_list(npc, category.prompt, options)
    if sel == nil then
        return nil
    end
    return category.recipes[sel]
end

function M.run(me, npc, categories, text)
    local category_options = {}
    for i, category in ipairs(categories) do
        category_options[i] = category.text
    end
    local category_sel = me:dialog_list(npc, text.greet, category_options)
    if category_sel == nil then
        return
    end
    local category = categories[category_sel]

    local recipe = choose_recipe(me, npc, category)
    if recipe == nil then
        return
    end

    local qty = 1
    if not category.equip then
        qty = me:dialog_number(npc, text.ask_qty(recipe.item), 1, 1, 100)
        if qty == nil then
            return
        end
    end

    local prompt = text.confirm(recipe.item, qty)
    local cost_items = {}
    for _, mat in ipairs(recipe.mats) do
        prompt = prompt .. "\r\n#i" .. mat[1] .. "# #t" .. mat[1] .. "# " .. mat[2] * qty .. text.unit
        cost_items[mat[1]] = mat[2] * qty
    end
    local meso = recipe.cost * qty
    if meso > 0 then
        prompt = prompt .. "\r\n#i4031138# " .. meso .. " 메소"
    end
    if not me:dialog_yes_no(npc, prompt) then
        return
    end

    if me:meso() < meso then
        me:dialog(npc, text.no_meso)
        return
    end
    local cost = { item = cost_items }
    if meso > 0 then
        cost.meso = meso
    end
    local code = me:exchange(cost, { item = { [recipe.item] = (recipe.gain or 1) * qty } })
    if code ~= ExchangeResult.OK then
        me:dialog(npc, text.failed)
        return
    end
    me:dialog(npc, text.done, false, text.done_next == true)
end

return M

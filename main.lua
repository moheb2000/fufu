local main_font = font("main", "./assets/IMFellEnglish-Regular.ttf")

narrate("What's your name? Can you tell me?", { text_color = "#2244ff", font = main_font, font_size = 22 })

local result = choice({ "Alex", "John", "Edward" }, { text_color = "#44ff22", font = main_font, font_size = 14 })
local name = ""
if result == 1 then
  name = "Alex"
elseif result == 2 then
  name = "John"
else
  name = "Edward"
end

local me = character(name, "#ff0000")

say(me, "Hello, my name is " .. me.name)
say(me, "I really like this visual novel engine", {
  color = "#55ffff",
  font = main_font
})

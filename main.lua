print(get_engine_version())
print(get_game_version())

play_music("assets/music2.mp3", true)

local main_font = font("main", "./assets/IMFellEnglish-Regular.ttf")

splash("assets/icon.png", "#222222", 2000)

bg("assets/background.png", { fade = true, originx = "center", originy = "center" })

narrate("What's your name? Can you tell me?", { text_color = "#2244ff", font = main_font, font_size = 22 })

pause_music()

local result = choice({ "Alex", "John", "Edward" }, { text_color = "#44ff22", font = main_font, font_size = 14 })
local name = ""
if result == 1 then
  name = "Alex"
elseif result == 2 then
  name = "John"
else
  name = "Edward"
end

local sound = "./assets/sound.wav"

local me = character(name, "#ff0000")

resume_music()

play_sound(sound)

say(me, "Hello, my name is " .. me.name)

stop_music()

play_sound(sound)

say(me, "I really like this visual novel engine", {
  color = "#55ffff",
  font = main_font
})

play_sound(sound)

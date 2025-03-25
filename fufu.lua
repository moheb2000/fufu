---@class font
---@field name string
---@field path string

---@param name string The font's name
---@param path string The font's path in filesystem
---@return font font A table containing the font's data
function font(name, path) end

---@class character
---@field name string
---@field color string

---@param name string The character's name
---@param color string The color used to show the character's name
---@return character character A table containing the caracter's data
function character(name, color) end

---@class properties
---@field font font?
---@field color string?
---@field font_size number?

---@param text string The text said by narrator
---@param properties properties? A table containing properties of the text
function narrate(text, properties) end

---@param character character The character that says the dialog
---@param text string The dialog said by character
---@param properties properties? A table containing properties of the text
function say(character, text, properties) end

---@param options string[]
---@param properties properties? A table containing properties of the text options
---@return result number The result of what user chose
function choice(options, properties) end

---@param path string The path to music for playing
---@param loop boolean? whether the music should loop or not
function play_music(path, loop) end

function stop_music() end

function pause_music() end

function resume_music() end

---@param path string the path to the sound for playing
function play_sound(path) end

---@return version string the engine version
function get_engine_version() end

---@return version string the game version
function get_game_version() end

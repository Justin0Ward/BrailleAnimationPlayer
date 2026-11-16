
-- Aseprite: Export animation frames as binary text files
-- 0 = white, 1 = black

local sprite = app.activeSprite

if not sprite then
  app.alert("No active sprite!")
  return
end

-- Require the sprite to be saved first.
if sprite.filename == "" then
  app.alert("Please save your sprite first.")
  return
end

local baseDir = app.fs.filePath(sprite.filename)
local baseName = app.fs.fileTitle(sprite.filename)

local outputFolder = app.fs.joinPath(baseDir, baseName)
app.fs.makeDirectory(outputFolder)

local pc = app.pixelColor

-- Make a safe folder name from a layer name.
local function safeName(name)
  return name:gsub("[^%w%-_]", "_")
end

-- Export one layer across all animation frames.
local function exportLayer(layer)

  local layerFolder = app.fs.joinPath(
    outputFolder,
    safeName(layer.name)
  )

  app.fs.makeDirectory(layerFolder)

  for frameNum = 1, #sprite.frames do

    -- Create a full-size white canvas.
    local img = Image(
      sprite.width,
      sprite.height,
      ColorMode.RGB
    )

    img:clear(nil, pc.rgba(255, 255, 255, 255))

    -- Get this layer's cel for the current frame.
    local cel = layer:cel(frameNum)

    if cel then
      -- Draw the cel at its correct canvas position.
      img:drawImage(
        cel.image,
        cel.position,
        cel.opacity
      )
    end

    -- Convert pixels into text.
    local lines = {}

    for y = 0, sprite.height - 1 do

      local row = {}

      for x = 0, sprite.width - 1 do

        local pixel = img:getPixel(x, y)

        local r = pc.rgbaR(pixel)
        local g = pc.rgbaG(pixel)
        local b = pc.rgbaB(pixel)

        -- Convert brightness to black or white.
        local brightness = (r + g + b) / 3

        if brightness < 128 then
          row[#row + 1] = "1"
        else
          row[#row + 1] = "0"
        end

      end

      lines[#lines + 1] = table.concat(row)
    end

    -- Save this frame.
    local filename = string.format(
      "frame_%03d.txt",
      frameNum
    )

    local filepath = app.fs.joinPath(
      layerFolder,
      filename
    )

    -- Ensure the output directory exists.
    app.fs.makeAllDirectories(layerFolder)

    local file, err = io.open(filepath, "w")

    if not file then
      app.alert(
        "Could not write file!\n\n" ..
        "Path:\n" .. filepath ..
        "\n\nError:\n" .. tostring(err)
      )
      return false
    end

    file:write(sprite.width, " ", sprite.height, "\n")
    file:write(table.concat(lines, "\n"))
    file:close()

  end

  return true
end

-- Recursively find image layers, including those in groups.
local function exportLayers(layers)

  for _, layer in ipairs(layers) do

    if layer.isGroup then
      if not exportLayers(layer.layers) then
        return false
      end

    elseif layer.isImage then
      if not exportLayer(layer) then
        return false
      end
    end

  end

  return true
end

if exportLayers(sprite.layers) then
  app.alert(
    "Export complete!\n\nSaved to:\n" .. outputFolder
  )
end

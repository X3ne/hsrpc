Get-ChildItem -Filter 'Character_*_Icon.webp' | ForEach-Object {
    if ($_ -match 'Character_(.+?)_Icon\.webp') {
        $characterName = $matches[1].ToLower()
        $outputFile = "char_$characterName.png"
        magick $_.FullName -resize 512x512 "$outputFile"
    }
}

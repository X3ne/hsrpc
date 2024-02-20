<p align="center"><img src="https://socialify.git.ci/X3ne/hsrpc/image?description=1&font=Inter&language=1&logo=https%3A%2F%2Fgithub.com%2FX3ne%2Fhsrpc%2Fblob%2Fmain%2Fassets%2Ficon.png%3Fraw%3Dtrue&name=1&owner=1&stargazers=1&theme=Auto" alt="project-image"></p>

<p id="description">This project is still in the very early stages of development. It still lacks many features and probably has a few bugs and performance problems.</p>

<h2>🖼️ Project Screenshots:</h2>

<div align="center">
  <img src="https://cdn.discordapp.com/attachments/568052462716583948/1208268522627670046/CaXsVbp.png?ex=65e2aac0&is=65d035c0&hm=fcddb9f18578452c0a036b4e231bcf2315f1f0a475c7dd2f10db92c72c432c98&" alt="project-screenshot" width="300"/>
</div>

<h2>💻 Platforms:</h2>

- [x] Windows
- [ ] Linux
- [ ] MacOs

<h2>🗣️ Languages:</h2>

- [x] English
- [ ] French

<h2>🤓 Technical information:</h2>

This program works with [tesseract](https://github.com/tesseract-ocr/tesseract), an open source text recognition software. hrpc extracts no data from the game and is based solely on what tesseract recognizes on the window (so some results may be wrong)

<h2>🛠️ Installation Steps:</h2>

<p>1. Install tesseract-ocr</p>

[tesseract github](https://github.com/tesseract-ocr/tesseract?tab=readme-ov-file#installing-tesseract)

or with choco

```
choco install tesseract
```

<p>2. Download hsrpc.exe</p>

Download the file named `hsrpc_windows_amd64.zip` and extract the executable (.exe) to the desired location

[Latest release](https://github.com/X3ne/hsrpc/releases/latest)

<p>3. Launch</p>

Just run the executable and off you go

<h2>🏗️ Build Steps:</h2>

<p>1. Clone this project</p>

```
git clone https://github.com/X3ne/hsrpc
```

<p>2. Install go-winres</p>

```
go install github.com/tc-hib/go-winres@latest
```

<p>3. Run go generate (used to bundle static files and set executable icon)</p>

```
go generate ./...
```

<p>4. Build</p>

```
go build -ldflags="-X github.com/X3ne/hsrpc/src/consts.Version=X.X.X -H=windowsgui"
```

<h2>🪲 Known issues:</h2>

- [ ] Sometimes, the position of the selected character is not the right one (especially when the background is too bright, e.g. on Jarilo-VI when there's snow in the background)
- [ ] I don't have these characters, but `Dan Heng Imbibitor Lunae` and `Topaz and Numby` seem sus names for ocr detection.

Menus:
- [ ] Data bank tab is not detected
- [ ] Achievements tab is not detected

- [ ] Fix gosec errors

<h2>⚒️ Improvements:</h2>

- [x] Add support for more game resolution (tested with fullscreen 1920x1080, 2560x1080 & windowed 1768x992)
- [x] Add more game menus status
- [ ] Clean some code
- [x] Create scripts to create data csv
- [x] Build to .exe
- [ ] Add support for Simulated Universe
- [ ] Add support for Calyx
- [ ] Add support for Cavern of corrosion
- [ ] Add support for Echo of war
- [ ] Add support for cut scenes
- [ ] Add support for the Trailblazer
- [ ] Remove the tesseract install step (maybe try to use [GetText](https://pkg.go.dev/github.com/go-vgo/robotgo#GetText) from robotgo)
- [ ] I want to add more infos for the selected character in the character tab (like the character name, level...)
- [x] Auto calculate coords based on app resolution
- [ ] Reimplement the personalization of coordinates in GUI

GUI app:
- [x] Add persistent state to GUI configuration
- [ ] Add field validation


<h2>🎨 Credits:</h2>

The assets and data for the discord presence come from the [Honkai Star Rail wiki](https://honkai-star-rail.fandom.com/wiki/Honkai:_Star_Rail_Wiki)

The [app icon](https://www.deviantart.com/mhesagnta/art/Chibi-Silver-Wolf-Honkai-StarRail-Render-965316702) by mhesagnta

Image assets are intellectual property of HoYoverse, © All rights reserved by miHoYo

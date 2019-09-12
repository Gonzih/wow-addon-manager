# WoW Addon Updater [![CircleCI](https://circleci.com/gh/Gonzih/wow-addon-manager.svg?style=svg)](https://circleci.com/gh/Gonzih/wow-addon-manager)

Simple go program that synchronizes local World of Warcraft addons from curseforge.

Works on Windows, Linux and Mac!

## Installation

```
go get github.com/Gonzih/wow-addon-manager
```

## Updating

```
go get -u github.com/Gonzih/wow-addon-manager
```

## Usage

Create `addons.yaml` file your `PATH_TO_WOW/Interface/Addons` folder with list of addons

```addons.yaml
addons:
- bartender4
- atlas
- atlas-classicwow
- classicauradurations
- classiccastbars
- classiclfg
- deadly-boss-mods
- details
- itemtooltipprofessionicons
- mapfader
- mik-scrolling-battle-text
- omni-cc
- quartz
- questie
- shadowed-unit-frames
- vendor-price
- weakauras-2
- whats-training
- xtolevel-classic
```

Each addon in this list is just part of url from curseforge,
so for `https://www.curseforge.com/wow/addons/atlas` it would be `atlas`,
for `https://www.curseforge.com/wow/addons/atlas-classicwow/` it would be `atlas-classicwow`.

Run the program
```
wow-addon-manager --addons-dir PATH_TO_WOW/Interface/Addons
```

## Windows notes

You will have first install go and git on your system for this to work.

Just get msi files from those urls:
* https://git-scm.com/download/win
* https://golang.org/dl/

Now you can open powershell or command prompt and install this program:
```
go get github.com/Gonzih/wow-addon-manager
```

To simplify update process one can create a simple `update.bat` file that contains following:

For classic WoW:
```
C:\Users\myusername\go\bin\wow-addon-manager.exe -addons-dir "C:\Program Files (x86)\World of Warcraft\_classic_\Interface\AddOns"
```

For retail WoW:
```
C:\Users\myusername\go\bin\wow-addon-manager.exe -addons-dir "C:\Program Files (x86)\World of Warcraft\Interface\AddOns"
```

`addons.yaml` file goes in to `C:\Program Files (x86)\World of Warcraft\_classic_\Interface\AddOns` for classic or `C:\Program Files (x86)\World of Warcraft\Interface\AddOns` for retail.

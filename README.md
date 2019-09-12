# WoW Addon Updater [![CircleCI](https://circleci.com/gh/Gonzih/wow-addon-manager.svg?style=svg)](https://circleci.com/gh/Gonzih/wow-addon-manager)

Simple go program that synchronizes local World of Warcraft addons from curseforge.

## Installation

```
   go get github.com/Gonzih/wow-addon-manager
```

## Updating

```
   go get -u github.com/Gonzih/wow-addon-manager
```

## Usage

Create addons.yaml file your `PATH_TO_WOW/Interface/Addons` folder with list of addons

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


Run the program

```
   wow-addon-manager --addons-dir PATH_TO_WOW/Interface/Addons
```

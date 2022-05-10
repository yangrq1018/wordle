# A Wordle Game
The "wordle" game is inspired by the game on New York Times.
    

## Quick start
> If you don't have a Powershell installation on your machine (version 7+ recommended), please install
> one from [here](https://docs.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows?view=powershell-7.2). You can ignore this step if you are on unix- systems, where ANSI colored supported terminal is standard.

Run the following command in terminal

`./wordle.exe`

Type a five-letter word (case-insensitive) and press enter until you got the word or 
used up all chances (by default six).

By default the program terminates after one round of play. To continue playing after winning or losing

`./wordle.exe -loop`

To cheat by displaying the secret word before the first guess

`./wordle.exe -cheat`

## Using the Powershell scripts (.ps1)
To run the game more conveniently on Windows platform, it is recommended to add "run powershell" option
to the right-click context menu of file type `.ps1` (Powershell scripts). A regedit script is included in the release, run it by double clicking. It should add a "Run with Powershell 7" option when you right
click `game.ps1`.

## Customizing dictionary
A text dictionary file (lowercase, one word each line). Multiple dictionaries are provided
with the binary releases, under `dictionary/`. You are free to make your own, just remember
to set parameter `-dict` to point to your dictionary file.

## Using the binary release
A github action workflow has been setup to build a binary release upon the push of version tags
(e.g. v1.0.0). Please go to [Release](https://github.com/yangrq1018/wordle/releases) to download
the latest version. Simply download the zip file matching your operating system and unpack it anywhere.

## Known issues
- Could randomly pick a strange word as secret
- Could pick a tense-variant of verb as secret
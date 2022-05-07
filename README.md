# A Wordle Game
The "wordle" game is inspired by the game on New York Times.

Prerequisites (on Windows platform):  
    - A Powershell terminal (supports ANSI escape sequences) which can print colorfully  
    - A text dictionary file (lowercase, one word each line)

A powershell 7 setup program should be included in the distribution.

## Usage
Run the following command in a Powershell terminal

`wordle.exe`

Type a five-letter word (case-insensitive) and press enter until you got the word or 
used up all chances.

## Cheat

`./wordle -cheat`

to set the secret word beforehand.

## About running the Powershell scripts (.ps1)
Run the regedit script to register "Run with Powershell 7" item in right-click context menu.

## Pick frequent word (yet to be done)
Choose a frequently used word from the dictionary.


## Known issues
- Could randomly pick a strange word as secret
- Could pick a tense-variant of verb as secret
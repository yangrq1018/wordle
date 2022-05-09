if (Test-Path -Path build) {
    Remove-Item -Recurse -Path build
}
New-Item -ItemType "directory" -Path build
Copy-Item -Path scripts/game.ps1,scripts/*.reg -Destination build
Copy-Item -Recurse -Path dictionary build
go build -o ./build/wordle.exe ./main
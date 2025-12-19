go build -o devinspect.exe

Remove-Item -Path C:\Tools\devinspect.exe -Force -Recurse
move .\devinspect.exe C:\Tools\

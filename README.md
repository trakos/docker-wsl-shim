This is a very simple shim for docker.exe and docker-compose.exe that allows using docker in wsl without Docker Desktop.

Manually install docker in WSL2, make sure it works inside the linux system, and then you can use this commands in
windows as if docker was running with Docker for Desktop.

This command uses only default WSL instance. To change that, simply assign desired name to `defaultDistro`.

It should take care of path substitutions, e. g. changing
 - `\\wsl$\system-name\var\opt` to `/var/opt` (ony if system-name is the defaultDistro)
 - `C:/Users/username/Desktop` to `/mnt/c/Users/username/Desktop`. It cals `wsl wslpath` to do this


Build:
```
go.exe build -o out\docker.exe docker.go helper.go
go.exe build -o out\docker-compose.exe docker-compose.go helper.go
```
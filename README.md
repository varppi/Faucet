
<img width=500 src="https://github.com/user-attachments/assets/d5a8c5bb-3cbb-4463-9390-c5e1aaeefaa2"></img>

[![Go/Golang](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
[![Webview](https://img.shields.io/badge/Google%20Chrome-4285F4?style=for-the-badge&logo=GoogleChrome&logoColor=white)](https://img.shields.io/badge/Google%20Chrome-4285F4?style=for-the-badge&logo=GoogleChrome&logoColor=white)
[![CSS3](https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white)](https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white)
[![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white)](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white)
[![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E)](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E)
[![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
[![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
[![Windows 11](https://img.shields.io/badge/Windows%2011-%230079d5.svg?style=for-the-badge&logo=Windows%2011&logoColor=white)](https://img.shields.io/badge/Windows%2011-%230079d5.svg?style=for-the-badge&logo=Windows%2011&logoColor=white)

### Easy to use encrypted cross platform file stealer with an interactive GUI. Yeah that summarizes it pretty well.
btw. feel free to contribute, I know my JS, HTML and CSS aren't the cleanest and greatest :D

## <b>Server Installation</b>
### Precompiled:
#### Go to releases, download the executable and run it.
### Compile yourself:
1. 
```
git clone https://github.com/R00tendo/Faucet
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
* password: Password used to encrypt the traffif.
* listen: Server listener. Format: `<host>:<port>`.
* lootdir: Where Faucet will store downloaded files.
* autodownload: A file/directory to download the second the agent connects. NOTE: %user% will be translated to the home directory.

### I don't really want skids using this so you have to have enough knowledge to figure out the rest yourself kekw

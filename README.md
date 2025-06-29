# vcurl - A Simple TUI for Sending HTTP Requests

`vcurl` is a text-based user interface (TUI) tool designed to send HTTP requests easily from your terminal. It allows you to quickly test and interact with RESTful APIs and other HTTP services without needing to write complex commands.

## Features
- Send GET, POST, PUT, PATCH and DELETE requests.
- Set custom headers and body data.
- View response status, headers, and body.
- Lightweight and simple to use.

## Installation

### Linux

Run the `vcurl-linux-install.sh` script. Change the path in the script if you want to.

To install manually:

1. Download the latest release from the [GitHub Releases Page](https://github.com/jiura/vcurl/releases)
2. Extract the `vcurl-linux-vX.Y.Z.tar.gz` file
3. Move the `vcurl` binary to `/usr/local/bin` (or your path of choosing):
```bash
sudo mv vcurl /usr/local/bin/vcurl
```
4. Make sure it's executable:
```bash
sudo chmod +x /usr/local/bin/vcurl
```

### Windows

Run the `vcurl-windows-install.bat` script. Change the path in the script if you want to.

To install manually:

1. Download the latest release from the [GitHub Releases Page](https://github.com/jiura/vcurl/releases)
2. Extract the `vcurl-windows-vX.Y.Z.zip` file to a folder, such as `C:\Program Files\vcurl`
3. Add the folder to your system's `PATH` to run it from any command prompt:
```bash
setx PATH "%PATH%;C:\Program Files\vcurl"
```

## Usage

Run:
```bash
vcurl
```

## Building

Make sure you have Go installed.

### Linux

```bash
go build -o path/you/want/vcurl
```

### Windows

```bash
GOOS=windows GOARCH=amd64 go build -o path/you/want/vcurl.exe
```

# scattered-identities

Google Drive and Sheets metadata parser.

## Usage

### Go
```bash
go run main.go -id 1ZXB6pnQEvym9pkDDoSxl49mxiiMAkzx0tS50OuwRt3Q
```

### Rust
```bash
rustc main.rs -o solve && ./solve 1ZXB6pnQEvym9pkDDoSxl49mxiiMAkzx0tS50OuwRt3Q
```

## Output

```console
$ go run main.go -id 1ZXB6pnQEvym9pkDDoSxl49mxiiMAkzx0tS50OuwRt3Q
[*] Target Document: https://drive.google.com/file/d/1ZXB6pnQEvym9pkDDoSxl49mxiiMAkzx0tS50OuwRt3Q/view
[+] Title: Nothing important
[+] Recovered Owner Email: cloveeen@gmail.com
[+] Extracted Username:   cloveeen
[*] Checking linked profile: https://www.youtube.com/@cloveeen
[!] Discovered Linked Real Name: Zuzanna Kasperek
[!] Linked Channel / Gaia ID:   UCEYNukARHDEMO8FtgRZGjXw
```

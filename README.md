# Welcome to Simul

Simul is a new work in progress, open source game with the primary goal of emulating brain lord

### Getting started:
- First see the install page from the ebiten website and follow the instructions that pertain to your system to install the necessary dependencies: [ebiten](https://ebiten.org/documents/install.html)
- Install the go compiler for your system: [Go](https://go.dev/dl/)
- `cd simulcast_game` then `go generate ./pkg/...` and then `go run simul.go`

### Project goals:

- To make a game similar in function to classics like Brain Lord, or LoZ LTTP
- To work wherever Go and Ebiten work, especially on non-windows non-x86 machines
  - The project targets Linux, MacOS, Windows
  - Raspberry Pi or Pi like boards, Pinebook, of course x86 type machines

### The necessary tools:
- Tiled
- Aseprite
- A text editor
- go install
  - golang.org/x/tools/cmd/stringer@latest
  - github.com/Racinettee/simulcast_game/cmd/collision_builder@latest
- others to potentially come...?

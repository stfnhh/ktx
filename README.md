<pre>
██╗  ██╗████████╗██╗  ██╗  
██║ ██╔╝╚══██╔══╝╚██╗██╔╝  
█████╔╝    ██║    ╚███╔╝   
██╔═██╗    ██║    ██╔██╗   
██║  ██╗   ██║   ██╔╝ ██╗  
╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝
</pre>

`ktx` is a simple terminal UI for quickly switching Kubernetes contexts.

It reads available contexts from `kubectl`, lets you filter and pick one, then runs `kubectl config use-context` for the selected item.

## Requirements

- Go 1.25+
- `kubectl` installed and configured
- At least one Kubernetes context in your kubeconfig

## Install

Build locally:

```bash
go build -o ktx .
```

Install to your Go bin path:

```bash
go install .
```

## Usage

Run:

```bash
ktx
```

If you built locally and did not move the binary into your `PATH`, run:

```bash
./ktx
```

## Controls

- `↑` / `↓` or `j` / `k`: move through contexts
- Type: filter contexts by name
- `Backspace`: clear filter characters
- `Enter`: select and switch context
- `q`, `esc`, or `ctrl+c`: quit without switching

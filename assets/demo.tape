# Regenerate with: vhs assets/demo.tape (from the repository root)

Output assets/demo.gif

Set Shell "bash"
Set FontSize 22
Set Width 1440
Set Height 760
Set Padding 24
Set Theme "Catppuccin Mocha"
Set TypingSpeed 45ms

Hide
Type `export D=/tmp/my-infra && rm -rf "$D" && mkdir -p "$D" && go build -o "$D/bin/tfpacker" . && cp assets/demo/main.tf assets/demo/tfpacker.config.yaml "$D" && cd "$D" && export PATH="$D/bin:$PATH" PS1='\[\e[1;35m\]$ \[\e[0m\]' && clear`
Enter
Sleep 4s
Show

Type "grep -E '^[a-z]' main.tf"
Sleep 400ms
Enter
Sleep 3s

Type "clear"
Enter
Type "cat tfpacker.config.yaml"
Sleep 400ms
Enter
Sleep 3s

Type "tfpacker --output tf"
Sleep 400ms
Enter
Sleep 3s

Type "clear"
Enter
Type "ls tf && cat tf/database.tf"
Sleep 400ms
Enter
Sleep 4s

# VHS Demo for Stamp
# Run: vhs demo.tape

Output demo.gif

Set FontSize 16
Set Width 1000
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20

# Setup
Hide
Type "go build -o stamp ./cmd/stamp"
Enter
Sleep 2s
Type "export DEMO_DIR=$(mktemp -d) && cp stamp $DEMO_DIR/ && cd $DEMO_DIR"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 500ms
Show

# Title
Type "# Stamp - A beautiful CLI for managing ADRs"
Sleep 1s
Enter
Sleep 500ms

# Initialize
Type "./stamp init"
Enter
Sleep 2s

# Create first ADR
Type "./stamp new 'Use PostgreSQL for persistence'"
Enter
Sleep 1.5s

# Create second ADR
Type "./stamp new 'Add Redis caching layer'"
Enter
Sleep 1.5s

# List ADRs
Type "./stamp list"
Enter
Sleep 3s

# Update status
Type "./stamp status 2 proposed"
Enter
Sleep 1.5s

Type "./stamp status 3 accepted"
Enter
Sleep 1.5s

# Show updated list
Type "./stamp list"
Enter
Sleep 3s

# Create superseding ADR
Type "./stamp new 'Switch to SQLite for simplicity'"
Enter
Sleep 1.5s

# Link ADRs
Type "./stamp link 4 2 supersedes"
Enter
Sleep 2s

# Show an ADR
Type "./stamp show 4"
Enter
Sleep 4s

# Final list
Type "./stamp list"
Enter
Sleep 3s

# End
Type ""
Sleep 2s

# Cleanup
Hide
Type "cd - && rm -rf $DEMO_DIR && rm stamp"
Enter

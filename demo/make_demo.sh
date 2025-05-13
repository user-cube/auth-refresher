#!/bin/bash
# Run demo_script and record
asciinema rec auth-refresher.cast -c "./demo_script.sh" --overwrite

# Convert to GIF
asciicast2gif auth-refresher.cast auth-refresher.gif

# Optimize GIF
gifsicle -O3 --colors 256 auth-refresher.gif -o auth-refresher-demo.gif

echo "✅ Done! Your demo GIF is ready as auth-refresher-demo.gif"
#!/bin/bash
# Run demo_script and record
asciinema rec $PROJECT_NAME.cast -c "./demo_script.sh" --overwrite

# Convert to GIF
asciicast2gif $PROJECT_NAME.cast $PROJECT_NAME.gif

# Optimize GIF
gifsicle -O3 --colors 256 $PROJECT_NAME.gif -o $PROJECT_NAME-demo.gif

echo "✅ Done! Your demo GIF is ready as $PROJECT_NAME-demo.gif"
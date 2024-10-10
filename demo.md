# Test

~~test~~
**bold**
*italic*

```bash
#!/bin/bash

# comment to make this file scroll
# comment to make this file scroll
# comment to make this file scroll
# comment to make this file scroll

set -u -o pipefail -u

# Array of ASCII art
ascii_art=(
"
  __
 /  \\
|  o |
 \\__/
  |||
"
"
   __
  (oo)
  (__)
  |  |
 //||\\\\
"
"
    _____
   /     \\
  | () () |
   \\  ^  /
    |||||
    |||||
"
"
 (\\(\\
 (-.-)
 o_(\")(\")
"
)

# Get a random number between 0 and the length of the array
random_index=$((RANDOM % ${#ascii_art[@]}))

# Print the random ASCII art
echo "${ascii_art[$random_index]}"
```

# Emoji Codes for Mac and iPad

A quick & dirty hack to generate an [importable text replacements propery list (.plist) file][1] to add system-wide support for emoji codes to Mac and iPad. 

I regret nothing. 

Learn more at https://sheesh.blog/system-wide-emoji-codes

`:v:` 

_NOTE: input data (`emoji.json`) originates from the [GitHub Emojis API][2]._ 

## Try it 

1. Clone this repository 

   ```
   git clone https://github.com/calebhailey/emojicodes.git
   ```

2. Compile

   ```
   cd emojicodes
   go build
   ```

3. Generate

   ```
   ./emojicodes > emojicodes.plist
   ```

[1]: https://support.apple.com/guide/mac-help/back-up-and-share-text-replacements-on-mac-mchl2a7bd795/mac
[2]: https://docs.github.com/en/rest/emojis
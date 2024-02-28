# GUI for FFmpeg

###### [Russian](README_ru.md) | English

A simple frontend for FFmpeg written in Go using Fyne. Originally made by [kor-elf](https://git.kor-elf.net/kor-elf/), on their private Gitea ([repo](https://git.kor-elf.net/kor-elf/gui-for-ffmpeg/)). This fork plans to add more features to the program, and more languages.

FFmpeg is a trademark of <a href="http://bellard.org/" target="_blank">Fabrice Bellard</a>, the creator of <a href="https://ffmpeg.org/about.html" target="_blank">FFmpeg</a>.

Both this fork and the original project is licensed under MIT (see [LICENSE](LICENSE)) and uses third party libraries which are distributed under their own terms (see [LICENSE-3RD-PARTY.txt](LICENSE-3RD-PARTY.txt)).

![images/screenshot-gui-for-ffmpeg.png](images/screenshot-gui-for-ffmpeg.png) 

To download, check the releases here: https://github.com/lostdusty/gui-for-ffmpeg/releases/latest

## Install using go and fyne:
Run these commands on your terminal:
1. `go install fyne.io/fyne/v2/cmd/fyne@latest`
2. ``fyne get git.kor-elf.net/kor-elf/gui-for-ffmpeg``

## Compile the project the manual way
Run these commands on your terminal:
1. ``git clone https://git.kor-elf.net/kor-elf/gui-for-ffmpeg.git``
2. ``cd gui-for-ffmpeg``
3. If you don't have Fyne or Go installed, follow this guide: https://docs.fyne.io/started/
4. After installing everything needed, simply run ``go run main.go``
5. To compile to other platforms, use ``go install github.com/fyne-io/fyne-cross@latest``
   * You MUST have Docker installed
   * Read more about fyne-cross here: https://github.com/fyne-io/fyne-cross
6. To cross compile to Windows, use: ``fyne-cross windows --icon icon.png --app-id "." -name "gui-for-ffmpeg"``, or to cross compile to Linux, run: ``fyne-cross linux --icon icon.png --app-id "." -name "gui-for-ffmpeg"``
   * The output will be on the folder `fyne-cross/bin`
<!-- 8. В папку **fyne-cross/bin/linux-amd64** или **fyne-cross/bin/windows-amd64** копируете:
   * icon.png
   * data
   * languages
   * LICENSE
   * LICENSE-3RD-PARTY.txt
<p><strong>Структура должна получиться такая:</strong></p>
<img src="images/screenshot-folder-structure.png"> 
Will modify how this works later on the project-->

## Working with translations:
1. ``go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest``
2. ``goi18n merge -sourceLanguage ru -outdir languages languages/active.\*.toml languages/translate.\*.toml``
3. In the **languages/translate.*.toml** files replace the text into the desired language
4. ``goi18n merge -sourceLanguage ru -outdir languages languages/active.\*.toml languages/translate.\*.toml``

See more about translations here: https://github.com/nicksnyder/go-i18n
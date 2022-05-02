package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Plist struct {
	XMLName xml.Name `xml:"plist"`
	Settings Settings `xml:"array"`
	Version string `xml:"version,attr"`
}

type Settings struct {
	XMLName          xml.Name `xml:"array"`
	TextReplacements []TextReplacement
}

type TextReplacement struct {
	XMLName       xml.Name `xml:"dict"`
	PhraseLabel   string   `xml:"phrase-key",innerxml`
	PhraseString   string   `xml:"phrase-string",innerxml`
	ShortcutLabel string   `xml:"shortcut-key",innerxml`
	ShortcutString string   `xml:"shortcut-string",innerxml`
}

func (s Settings) AddTextReplacement(phrase string, shortcut string) Settings {
	phrase = EscapeUnicodeCodePoints(strings.Split(phrase,"-"))
	tr := TextReplacement{
		PhraseLabel:   "phrase",
		PhraseString:   phrase,
		ShortcutLabel: "shortcut",
		ShortcutString: shortcut,
	}
	s.TextReplacements = append(s.TextReplacements, tr)
	return s
}

func GetSortedKeysFromMap(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range(m) {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func EscapeUnicodeCodePoints(codepoints []string) string {
	escaped := ""
	for _,c := range(codepoints) {
		// escaped = escaped + `\u` + c
		escaped = escaped + "&#x" + c + ";"
	}
	return escaped
}

func GetEmojis() (map[string]string, error) {
	data := make(map[string]string)
	resp, err := http.Get("https://api.github.com/emojis")
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		log.Fatalf("ERROR: %v %s (%s)\n", resp.StatusCode, http.StatusText(resp.StatusCode), "https://api.github.com/emojis")
		return nil, err
	} else {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err)
		} else {
			err = json.Unmarshal(b, &data)
			if err != nil {
				log.Fatalf("ERROR: %s\n", err)
			}
		}
	}
	return data, nil
}

func GenerateEmojiPlist(emojis map[string]string) (Settings, error) {
	emojicodes := Settings{}
	shortcodes := GetSortedKeysFromMap(emojis)
	for i := range(shortcodes) {
		shortcode := fmt.Sprintf(":%s:", shortcodes[i])
		phrase := emojis[shortcodes[i]]
		if strings.Index(phrase, "unicode") > 0 {
			phrase = strings.TrimPrefix(phrase, "https://github.githubassets.com/images/icons/emoji/unicode/")
			phrase = strings.TrimSuffix(phrase, ".png?v8")
			emojicodes = emojicodes.AddTextReplacement(phrase, shortcode)
		}
	}
	// for i,v := range emojis {
	// 	if (strings.Index(v, "unicode") < 0) {
	// 		fmt.Printf("%s: %s\n", i,v)
	// 	}
	// }
	// fmt.Printf("%v\n", emojicodes.TextReplacements)
	// fmt.Printf("Emojis: %v\n", len(emojis))
	return emojicodes, nil
}

func PrintEmojiPlist(emojis Settings) error {
	const plistheader = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n"
	plist := Plist{ Settings: emojis, Version: "1.0" }
	plistbytes, err := xml.MarshalIndent(plist, "", "  ")
	plistbytes = []byte(xml.Header + plistheader + string(plistbytes))
	pliststring := string(plistbytes)
	pliststring = strings.Replace(pliststring, "phrase-key", "key", -1)
	pliststring = strings.Replace(pliststring, "phrase-string", "string", -1)
	pliststring = strings.Replace(pliststring, "shortcut-key", "key", -1)
	pliststring = strings.Replace(pliststring, "shortcut-string", "string", -1)
	if err == nil {
		fmt.Printf("%s\n", pliststring)
	} else {
		fmt.Printf("ERROR: %s\n", err)
	}
	return nil
}

func main() {
	github_emojis, err := GetEmojis()
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	emoji_plist, err := GenerateEmojiPlist(github_emojis)
	PrintEmojiPlist(emoji_plist)
}

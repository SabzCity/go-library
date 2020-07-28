/* For license and copyright information please see LEGAL file in repository */

package www

import (
	"bytes"
	"regexp"
	"strconv"

	"../assets"
	"../json"
	"../log"
)

var htmlComment = regexp.MustCompile(`(<!--.*?-->)|(<!--[\w\W\n\s]+?-->)`)

// MinifyHTML use to minify html data.
func MinifyHTML(html []byte) []byte {
	var f = newLine.ReplaceAll(html, []byte{})
	f = htmlComment.ReplaceAll(f, []byte{})
	return f
}

// mixHTMLToJS add given HTML file to specific part of JS file.
func mixHTMLToJS(jsFile, htmlFile *assets.File) {
	var loc = bytes.Index(jsFile.Data, []byte("HTML: ("))
	if loc < 0 {
		log.Warn(htmlFile.FullName, "html file can't add to", jsFile.FullName, "file due to bad JS.")
		return
	}

	var graveAccentIndex int
	graveAccentIndex = bytes.IndexByte(jsFile.Data[loc:], '`')

	var minifiedHTML = MinifyHTML(htmlFile.Data)

	var jsFileData = make([]byte, 0, len(jsFile.Data)+len(minifiedHTML))
	jsFileData = append(jsFileData, jsFile.Data[:loc+graveAccentIndex+1]...)
	jsFileData = append(jsFileData, minifiedHTML...)
	jsFileData = append(jsFileData, jsFile.Data[loc+graveAccentIndex+1:]...)
	jsFile.Data = jsFileData
}

// mixHTMLTemplateToJS add given HTML template file to specific part of JS file.
func mixHTMLTemplateToJS(jsFile, htmlFile *assets.File, tempName string) {
	var loc = bytes.Index(jsFile.Data, []byte(tempName+`": (`))
	if loc < 0 {
		log.Warn(htmlFile.FullName, "html template file can't add to", jsFile.FullName, "file due to bad JS template.")
		return
	}

	var graveAccentIndex int
	graveAccentIndex = bytes.IndexByte(jsFile.Data[loc:], '`')

	var minifiedHTML = MinifyHTML(htmlFile.Data)

	var jsFileData = make([]byte, 0, len(jsFile.Data)+len(minifiedHTML))
	jsFileData = append(jsFileData, jsFile.Data[:loc+graveAccentIndex+1]...)
	jsFileData = append(jsFileData, minifiedHTML...)
	jsFileData = append(jsFileData, jsFile.Data[loc+graveAccentIndex+1:]...)
	jsFile.Data = jsFileData
}

// localizeHTMLFile make and returns number of localize file by number of language indicate in JSON localize text
func localizeHTMLFile(html *assets.File, lj localize) (files map[string]*assets.File) {
	files = make(map[string]*assets.File, len(lj))
	for lang, text := range lj {
		files[lang] = replaceLocalizeTextInHTML(html, text, lang)
	}
	return
}

func replaceLocalizeTextInHTML(html *assets.File, text []string, lang string) (newFile *assets.File) {
	newFile = html.Copy()
	newFile.Name += "-" + lang
	newFile.FullName = newFile.Name + "." + newFile.Extension
	newFile.Data = nil

	var htmlData = html.Data

	var replacerIndex int
	var bracketIndex int
	var textIndex uint64
	var err error
	for {
		replacerIndex = bytes.Index(htmlData, []byte("${LocaleText["))
		if replacerIndex < 0 {
			newFile.Data = append(newFile.Data, htmlData...)
			return
		}
		newFile.Data = append(newFile.Data, htmlData[:replacerIndex]...)
		htmlData = htmlData[replacerIndex:]

		bracketIndex = bytes.IndexByte(htmlData, ']')
		textIndex, err = strconv.ParseUint(string(htmlData[13:bracketIndex]), 10, 8)
		if err != nil {
			log.Warn("Index numbers in", html.FullName, "is not valid, double check your file for bad structures")
		} else {
			newFile.Data = append(newFile.Data, text[textIndex]...)
		}

		htmlData = htmlData[bracketIndex+2:]
	}
}

type localize map[string][]string

func (lj *localize) jsonDecoder(data []byte) (err error) {
	// TODO::: convert to generated code
	err = json.UnMarshal(data, lj)
	return
}
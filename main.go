// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// An example command-line tool that uses opennota/markdown to process markdown input.
package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/browser"

	"github.com/opennota/markdown"
)

var (
	allowhtml   bool
	tables      bool
	linkify     bool
	typographer bool
	xhtml       bool

	title          string
	rendererOutput string

	wg sync.WaitGroup
)

func readFromStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}

func readFromFile(fn string) ([]byte, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func readFromWeb(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func readInput(input string) ([]byte, error) {
	if input == "-" {
		return readFromStdin()
	}
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return readFromWeb(input)
	}
	return readFromFile(input)
}

func extractText(tok markdown.Token) string {
	switch tok := tok.(type) {
	case *markdown.Text:
		return tok.Content
	case *markdown.Inline:
		text := ""
		for _, tok := range tok.Children {
			text += extractText(tok)
		}
		return text
	}
	return ""
}

func writePreamble(w io.Writer) error {
	var opening string
	var ending string
	if xhtml {
		opening = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
  "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">`
		ending = " /"
	} else {
		opening = `<!DOCTYPE html>
<html>`
	}
	_, err := fmt.Fprintf(w, `%s
<head>
<meta charset="utf-8"%s>
<title>%s</title>
</head>
<body>
`, opening, ending, html.EscapeString(title))

	return err
}

func writePostamble(w io.Writer) error {
	_, err := fmt.Fprint(w, `</body>
</html>
`)
	return err
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer wg.Done()

	err := writePreamble(w)
	if err != nil {
		log.Println(err)
	}

	_, err = fmt.Fprint(w, rendererOutput)
	if err != nil {
		log.Println(err)
	}

	err = writePostamble(w)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(1)
}

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: mdtool [options] [inputfile|URL] [outputfile]

Options:
  +h[tml]         Enable HTML
  +l[inkify]      Enable autolinking
  +ta[bles]       Enable GFM tables
  +ty[pographer]  Enable typographic replacements
  +a[ll]          All of the above
  +x[html]        XHTML output

  -help           Display help

Use 'browser:' in place of the output file to get the output in a browser.
`)
	}
	flag.Parse()
	var documents []string
	for _, arg := range flag.Args() {
		switch arg {
		case "+html", "+h":
			allowhtml = true
		case "+linkify", "+l":
			linkify = true
		case "+tables", "+ta":
			tables = true
		case "+typographer", "+ty":
			typographer = true
		case "+t":
			fmt.Fprintf(os.Stderr, "ambiguous option: +t; did you mean +ta[bles] or +ty[pographer]?")
			os.Exit(1)
		case "+xhtml", "+x":
			xhtml = true
		case "+all", "+a":
			allowhtml = true
			linkify = true
			tables = true
			typographer = true
		default:
			documents = append(documents, arg)
		}
	}
	if len(documents) > 2 {
		flag.Usage()
		os.Exit(1)
	}
	if len(documents) == 0 {
		documents = []string{"-"}
	}

	data, err := readInput(documents[0])
	if err != nil {
		log.Fatal(err)
	}

	md := markdown.New(
		markdown.HTML(allowhtml),
		markdown.Tables(tables),
		markdown.Linkify(linkify),
		markdown.Typographer(typographer),
		markdown.XHTMLOutput(xhtml),
	)

	tokens := md.Parse(data)
	if len(tokens) > 0 {
		if heading, ok := tokens[0].(*markdown.HeadingOpen); ok {
			for i := 1; i < len(tokens); i++ {
				if tok, ok := tokens[i].(*markdown.HeadingClose); ok && tok.Lvl == heading.Lvl {
					break
				}
				title += extractText(tokens[i])
			}
			title = strings.TrimSpace(title)
		}
	}

	rendererOutput = md.RenderTokensToString(tokens)

	if len(documents) == 1 {
		writePreamble(os.Stdout)
		fmt.Println(rendererOutput)
		writePostamble(os.Stdout)
	} else if documents[1] == "browser:" {
		srv := httptest.NewServer(http.HandlerFunc(handler))
		wg.Add(1)
		err = browser.OpenURL(srv.URL)
		if err != nil {
			log.Fatal(err)
		}
		wg.Wait()
	} else {
		f, err := os.OpenFile(documents[1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Println(err)
			}
		}()

		err = writePreamble(f)
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.WriteString(rendererOutput)
		if err != nil {
			log.Fatal(err)
		}

		err = writePostamble(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

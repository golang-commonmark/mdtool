mdtool [![License](https://img.shields.io/badge/licence-BSD--2--Clause-blue.svg)](https://opensource.org/licenses/BSD-2-Clause)
======

An example command-line tool that uses [golang-commonmark/markdown](https://github.com/golang-commonmark/markdown) to process markdown input.

## Installation

    go get github.com/golang-commonmark/mdtool

## Usage

    $ mdtool -help
    Usage: mdtool [options] [inputfile|URL] [outputfile]
    
    Options:
      +h[tml]         Enable raw HTML
      +l[inkify]      Enable autolinking
      +ta[bles]       Enable GFM tables
      +ty[pographer]  Enable typographic replacements
      +a[ll]          All of the above
      +x[html]        XHTML output
    
      -help           Display this help
    
    Use 'browser:' in place of the output file to see the output in a browser.

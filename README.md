mdtool [![License](http://img.shields.io/:license-gpl3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.html)
======

An example command-line tool that uses [opennota/markdown](https://github.com/opennota/markdown) to process markdown input.

## Installation

    go get github.com/opennota/mdtool

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

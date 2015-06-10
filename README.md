mdtool
======

An example command-line tool that uses [opennota/markdown](https://github.com/opennota/markdown) to process markdown input.

## Installation

    go get github.com/opennota/mdtool

## Usage

    $ mdtool -help
    Usage: mdtool [options] [inputfile|URL] [outputfile]
    
    Options:
      +h[tml]         Enable HTML
      +l[inkify]      Enable autolinking
      +ta[bles]       Enable GFM tables
      +ty[pographer]  Enable typographic replacements
      +a[ll]          All of the above
      +x[html]        XHTML output
    
      -help           Display help
    
    Use 'browser:' in place of the output file to get the output in a browser.

## License

GNU GPL v3+

package display

import (
	"bytes"
	"fmt"
)

type Prefixer struct {
	// a prefix to pre-pend to new lines
	prefix string

	// the first time Parse is called, we need to print a prefix for the first
	// line. After printing, this should be set to false so subsequent parses
	// won't print the prefix unless a newline is found.
	firstLine bool

	// indicate if we're on a new line. If true, the next character should
	// be printed with the prefix provided.
	newLine bool

	// escape sequences can move the cursor around the lines. We need to adjust
	// for prefix width when this happens
	escapeSeq bool

	// create a buffer to store the current escape sequence
	escapeBuf string
}

// NewPrefixer returns a Prefixer struct
func NewPrefixer(prefix string) *Prefixer {
	return &Prefixer{
		prefix:    prefix,
		firstLine: true,
		newLine:   false,
		escapeSeq: false,
		escapeBuf: "",
	}
}

// Prefix will prefix lines with a prefix
func (p *Prefixer) Parse(data string) string {
	// create an empty buffer to apply transformations onto
	buffer := bytes.NewBufferString("")

	// if this is the first line, let's add the prefix
	if p.firstLine == true {
		buffer.WriteString(p.prefix)
		// unset firstLine
		p.firstLine = false
	}

	// range over the string to extract the runes
	for _, r := range data {
		// convert the rune into a character
		c := string(r)

		// as soon as we see an escape character
		if c == "\x1b" {
			p.escapeSeq = true
			p.escapeBuf += c
		} else if p.escapeSeq == true {
			// add character to escape seq buffer
			p.escapeBuf += c

			// check if escape sequence is complete
			if p.isSequenceEnd(c) {
				// check if goto or horizontal reset and modify
				if c == "G" && p.newLine == false {
					p.adjustHorizontalSeq()
				} else if c == "H" && p.newLine == false {
					p.adjustGotoSeq()
				}

				// add escape sequence to buffer
				buffer.WriteString(p.escapeBuf)

				// reset escape sequence state
				p.escapeSeq = false
				p.escapeBuf = ""
			}

		} else {
			// add the character (prefixed if this is a new line)
			if p.newLine == true {
				buffer.WriteString(fmt.Sprintf("%s%s", p.prefix, c))
			} else {
				buffer.WriteString(c)
			}

			// set new_line to true if the printed character was a newline
			if c == "\n" || c == "\r" {
				p.newLine = true
			} else {
				p.newLine = false
			}
		}

	}

	return buffer.String()
}

// isSequenceEnd looks for a character denoting the end of an escape sequence
func (p Prefixer) isSequenceEnd(char string) bool {

	terminators := "ABCDEFGHJKRSTfminsulhp"

	// range of the string to extract the runes
	for _, r := range terminators {
		// convert the rune to a character
		c := string(r)

		if char == c {
			return true
		}
	}

	return false
}

// adjustHorizontalSeq adjusts the horizontal reset to compensate for the prefix
func (p *Prefixer) adjustHorizontalSeq() {
	var number int
	// Scan the escape sequence to extract the current offset
	fmt.Sscanf(p.escapeBuf, "\x1b[%dG", number)
	// Reset the escape sequence to add the prefix width
	p.escapeBuf = fmt.Sprintf("\x1b[%dG", number+len(p.prefix))
}

// adjustGotoSeq adjusts an escape goto reset to compensate for the prefix
func (p *Prefixer) adjustGotoSeq() {
	var x, y int
	// Scan the escape sequence to extract the current coordinates
	fmt.Sscanf(p.escapeBuf, "\x1b[%d;%dH", y, x)
	// Reset the escape sequence to add the prefix width
	p.escapeBuf = fmt.Sprintf("\x1b[%d;%dH", y, x+len(p.prefix))
}

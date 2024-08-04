package extensions

import (
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/apricote/releaser-pleaser/internal/markdown/extensions/ast"
)

var sectionStartRegex = regexp.MustCompile(`^<!-- section-start (.+) -->`)
var sectionEndRegex = regexp.MustCompile(`^<!-- section-end (.+) -->`)

const (
	sectionTrigger     = "<!--"
	SectionStartFormat = "<!-- section-start %s -->"
	SectionEndFormat   = "<!-- section-end %s -->"
)

type sectionParser struct {
}

func (s *sectionParser) Open(_ gast.Node, reader text.Reader, _ parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()

	if result := sectionStartRegex.FindSubmatch(line); result != nil {
		reader.AdvanceLine()
		return ast.NewSection(string(result[1])), parser.HasChildren
	}

	return nil, parser.NoChildren
}

func (s *sectionParser) Continue(node gast.Node, reader text.Reader, _ parser.Context) parser.State {
	n := node.(*ast.Section)

	line, _ := reader.PeekLine()

	if result := sectionEndRegex.FindSubmatch(line); result != nil {
		if string(result[1]) == n.Name {
			reader.AdvanceLine()
			return parser.Close
		}
	}

	return parser.Continue | parser.HasChildren
}

func (s *sectionParser) Close(_ gast.Node, _ text.Reader, _ parser.Context) {
	// Nothing to do
}

func (s *sectionParser) CanInterruptParagraph() bool {
	return true
}

func (s *sectionParser) CanAcceptIndentedLine() bool {
	return false
}

var defaultSectionParser = &sectionParser{}

// NewSectionParser returns a new BlockParser that can parse
// a section block. Section blocks can be used to group various nodes under a parent ast node.
// This parser must take precedence over the parser.HTMLParser.
func NewSectionParser() parser.BlockParser {
	return defaultSectionParser
}

func (s *sectionParser) Trigger() []byte {
	return []byte(sectionTrigger)
}

type section struct {
}

// Section is an extension that allow you to use group content under a shared parent ast node.
var Section = &section{}

func (e *section) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewSectionParser(), 0),
	))
}
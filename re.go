/* Package go_re
 * @Author 砚池/Ivan
 * @Date 2024/04/10
 * @Description:
 */

package re

import (
	"github.com/samber/lo"
	"regexp"
	"strconv"
	"strings"
)

// PatternFlag info https://github.com/google/re2/wiki/Syntax
type PatternFlag string

const (
	FLAG_CASE_INSENSITIVE        PatternFlag = "i"
	FLAG_MULTILINE               PatternFlag = "m"
	FLAG_DOTALL                  PatternFlag = "s"
	FLAG_UNICODE_CHARACTER_CLASS PatternFlag = "U"
)

type Expression struct {
	regexp *regexp.Regexp
}

func (e *Expression) String() string {
	return e.regexp.String()
}

type Builder struct {
	prefix    *strings.Builder
	source    *strings.Builder
	suffix    *strings.Builder
	modifiers []PatternFlag
}

func ExpressionBuilder() *Builder {
	return &Builder{
		prefix:    &strings.Builder{},
		source:    &strings.Builder{},
		suffix:    &strings.Builder{},
		modifiers: []PatternFlag{},
	}
}

func (b *Builder) sanitize(pValue string) string {
	return strings.ReplaceAll(pValue, "[\\W]", "\\\\$0")
}

func (b *Builder) countOccurrencesOf(where, what string) int {
	return len(where) - len(strings.Replace(where, what, "", 1))/len(what)
}

func (b *Builder) Build() *Expression {
	sb := strings.Builder{}
	if len(b.modifiers) > 0 {
		sb.WriteString("(?")
		modifiers := strings.Join(lo.Map[PatternFlag, string](b.modifiers, func(item PatternFlag, _ int) string {
			return string(item)
		}), "")
		sb.WriteString(modifiers)
		sb.WriteString(")")
	}
	sb.WriteString(b.prefix.String())
	sb.WriteString(b.source.String())
	sb.WriteString(b.suffix.String())
	re, err := regexp.Compile(sb.String())
	if err != nil {
		panic(err)
	}
	return &Expression{regexp: re}
}

func (b *Builder) Add(pValue string) *Builder {
	b.source.WriteString(pValue)
	return b
}

func (b *Builder) AddBuilder(regex *Builder) *Builder {
	b.Group().Add(regex.Build().String()).EndGr()
	return b
}

func (b *Builder) StartOfLineWithPrefix(pEnable bool) *Builder {
	b.prefix.WriteString(lo.Ternary(pEnable, "^", ""))
	if !pEnable {
		nb := strings.Builder{}
		nb.WriteString(strings.Replace(b.prefix.String(), "^", "", 1))
		b.prefix = &nb
	}
	return b
}

func (b *Builder) StartOfLine() *Builder {
	return b.StartOfLineWithPrefix(true)
}

func (b *Builder) EndOfLineWithSuffix(sEnable bool) *Builder {
	b.suffix.WriteString(lo.Ternary(sEnable, "$", ""))
	if !sEnable {
		nb := strings.Builder{}
		nb.WriteString(strings.Replace(b.suffix.String(), "$", "", 1))
		b.suffix = &nb
	}
	return b
}

func (b *Builder) EndOfLine() *Builder {
	return b.EndOfLineWithSuffix(true)
}

func (b *Builder) Then(pValue string) *Builder {
	return b.Add("(?:" + b.sanitize(pValue) + ")")
}

func (b *Builder) Find(value string) *Builder {
	return b.Then(value)
}

// Maybe prefer one
func (b *Builder) Maybe(value string) *Builder {
	return b.Then(value).Add("?")
}

// MaybePreferZero prefer zero
func (b *Builder) MaybePreferZero(value string) *Builder {
	return b.Then(value).Add("??")
}

// ZeroOrOne prefer one
func (b *Builder) ZeroOrOne(value string) *Builder {
	return b.Maybe(value)
}

// ZeroOrOnePreferZero prefer one
func (b *Builder) ZeroOrOnePreferZero(value string) *Builder {
	return b.MaybePreferZero(value)
}

func (b *Builder) MaybeWithBuilder(regex *Builder) *Builder {
	return b.Group().AddBuilder(regex).EndGr().Add("?")
}

func (b *Builder) Anything() *Builder {
	return b.Add("(?:.*)")
}

func (b *Builder) AnythingBut(value string) *Builder {
	return b.Add("(?:[^" + b.sanitize(value) + "])*")
}

func (b *Builder) Dot() *Builder {
	return b.Add("(?:.)")
}

func (b *Builder) Something() *Builder {
	return b.Add("(?:.+)")
}

func (b *Builder) SomethingBut(value string) *Builder {
	return b.Add("(?:[^" + b.sanitize(value) + "])+")
}

func (b *Builder) LineBreak() *Builder {
	return b.Add("(?:\\n|(?:\\r\\n)|(?:\\r\\r))")
}

func (b *Builder) Br() *Builder {
	return b.LineBreak()
}

func (b *Builder) Tab() *Builder {
	return b.Add("(?:\\t)")
}

func (b *Builder) Word() *Builder {
	return b.Add("(?:\\w+)")
}

func (b *Builder) WordChar() *Builder {
	return b.Add("(?:\\w)")
}

func (b *Builder) NonWordChar() *Builder {
	return b.Add("(?:\\W)")
}

func (b *Builder) Digit() *Builder {
	return b.Add("(?:\\d)")
}

func (b *Builder) NonDigit() *Builder {
	return b.Add("(?:\\D)")
}

func (b *Builder) Space() *Builder {
	return b.Add("(?:\\s)")
}

func (b *Builder) NonSpace() *Builder {
	return b.Add("(?:\\S)")
}

func (b *Builder) WordBoundary() *Builder {
	return b.Add("(?:\\b)")
}

func (b *Builder) Any(value string) *Builder {
	return b.AnyOf(value)
}

func (b *Builder) AnyOf(value string) *Builder {
	return b.Add("[" + b.sanitize(value) + "]")
}

func (b *Builder) Range(args ...string) *Builder {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i := 1; i < len(args); i = i + 2 {
		from := b.sanitize(args[i-1])
		to := b.sanitize(args[i])
		sb.WriteString(from)
		sb.WriteString("-")
		sb.WriteString(to)
	}
	sb.WriteString("]")
	return b.Add(sb.String())
}

func (b *Builder) AddModifier(modifier PatternFlag) *Builder {
	b.modifiers = append(b.modifiers, modifier)
	return b
}

func (b *Builder) RemoveModifier(modifier PatternFlag) *Builder {
	b.modifiers = lo.DropWhile(b.modifiers, func(item PatternFlag) bool {
		return item == modifier
	})
	return b
}

func (b *Builder) WithAnyCaseEnable(enable bool) *Builder {
	if enable {
		b.AddModifier(FLAG_CASE_INSENSITIVE)
	} else {
		b.RemoveModifier(FLAG_CASE_INSENSITIVE)
	}
	return b
}

func (b *Builder) WithAnyCase() *Builder {
	return b.WithAnyCaseEnable(true)
}

func (b *Builder) SearchMultiLineEnable(enable bool) *Builder {
	if enable {
		b.AddModifier(FLAG_MULTILINE)
	} else {
		b.RemoveModifier(FLAG_MULTILINE)
	}
	return b
}

func (b *Builder) SearchMultiLine() *Builder {
	return b.SearchMultiLineEnable(true)
}

func (b *Builder) Multiple(value string, count ...int) *Builder {
	if count == nil {
		return b.Then(value).OneOrMore()
	}
	switch len(count) {
	case 1:
		return b.Then(value).Count(count[0])
	case 2:
		return b.Then(value).CountBetween(count[0], count[1])
	default:
		return b.Then(value).OneOrMore()
	}
}

// OneOrMore prefer more
func (b *Builder) OneOrMore() *Builder {
	return b.Add("+")
}

// OneOrMorePreferFewer prefer fewer
func (b *Builder) OneOrMorePreferFewer() *Builder {
	return b.Add("+?")
}

// ZeroOrMore prefer more
func (b *Builder) ZeroOrMore() *Builder {
	return b.Add("*")
}

// ZeroOrMorePreferFewer prefer fewer
func (b *Builder) ZeroOrMorePreferFewer() *Builder {
	return b.Add("*?")
}

func (b *Builder) Count(count int) *Builder {
	b.source.WriteString("{" + strconv.Itoa(count) + "}")
	return b
}

// CountBetween prefer more
func (b *Builder) CountBetween(from, to int) *Builder {
	b.source.WriteString("{" + strconv.Itoa(from) + "," + strconv.Itoa(to) + "}")
	return b
}

// CountBetweenPreferFewer prefer more
func (b *Builder) CountBetweenPreferFewer(from, to int) *Builder {
	b.source.WriteString("{" + strconv.Itoa(from) + "," + strconv.Itoa(to) + "}?")
	return b
}

func (b *Builder) AtLeast(from int) *Builder {
	return b.Add("{").Add(strconv.Itoa(from)).Add(",}")
}

func (b *Builder) Or(value string) *Builder {
	b.prefix.WriteString("(?:")
	opened := b.countOccurrencesOf(b.prefix.String(), "(")
	closed := b.countOccurrencesOf(b.prefix.String(), ")")
	if opened >= closed {
		nb := strings.Builder{}
		nb.WriteString(")" + b.suffix.String())
		b.suffix = &nb
	}
	b.Add(")|(?:")
	if !lo.IsNil(value) {
		b.Then(value)
	}
	return b
}

func (b *Builder) OneOf(values ...string) *Builder {
	if !lo.IsNil(values) && len(values) > 0 {
		b.Add("(?:")
		for i := 0; i < len(values); i++ {
			b.Add("(?:").Add(values[i]).Add(")")
			if i < len(values)-1 {
				b.Add("|")
			}
		}
		b.Add(")")
	}
	return b
}

func (b *Builder) Capture() *Builder {
	return b.CaptureWithName("")
}

func (b *Builder) CaptureWithName(name string) *Builder {
	b.suffix.WriteString(")")
	if len(name) > 0 {
		return b.Add("(?P<" + name + ">")
	} else {
		return b.Add("(")
	}
}

func (b *Builder) Capt() *Builder {
	return b.Capture()
}

func (b *Builder) CaptWithName(name string) *Builder {
	return b.CaptureWithName(name)
}

func (b *Builder) Group() *Builder {
	b.suffix.WriteString(")")
	return b.Add("(?:")
}

func (b *Builder) EndCapture() *Builder {
	if strings.Index(b.suffix.String(), ")") != -1 {
		b.suffix = SetLength(b.suffix, b.suffix.Len()-1)
		return b.Add(")")
	} else {
		panic("Can't end capture (group) when it not started")
	}
}

func SetLength(sb *strings.Builder, length int) *strings.Builder {
	nb := strings.Builder{}
	nb.WriteString(sb.String()[:length])
	return &nb
}

func (b *Builder) EndCapt() *Builder {
	return b.EndCapture()
}

func (b *Builder) EndGr() *Builder {
	return b.EndCapture()
}

func (e *Expression) Test(str string) bool {
	var ret = false
	if !lo.IsNil(str) {
		ret = e.regexp.MatchString(str)
	}
	return ret
}

func (e *Expression) GetText(str, group string) []string {
	matches := e.regexp.FindAllStringSubmatch(str, -1)
	groupNames := e.regexp.SubexpNames()
	var res []string
	for _, match := range matches {
		for index, name := range groupNames {
			if name == group {
				res = append(res, strings.TrimSpace(match[index]))
			}
		}
	}
	return res
}

func (e *Expression) GetTextGroups(str string, group int) []string {
	panic("implement me")
}

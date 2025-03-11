# go-re

see more https://github.com/google/re2/wiki/Syntax

## Install

```bash
go get -u github.com/GoCarnival/go-re
```

## Usage demo

```go
func TestRegex(t *testing.T) {
	//目标字符串
	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	//pattern := `[0-9]+\.[0-9]+` //正则表达式
	expression := ExpressionBuilder().
		Digit().OneOrMore().
		Then("\\.").
		Digit().OneOrMore().
		Build()
	t.Log(expression.String())
	t.Log(expression.Test(searchIn))
	t.Log(expression.Regexp().ReplaceAllString(searchIn, "##.#"))
}
```

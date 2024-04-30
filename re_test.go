/* Package re
 * @Author 砚池/Ivan
 * @Date 2024/04/10
 * @Description:
 */

package re

import (
	"regexp"
	"testing"
)

func TestRe(t *testing.T) {
	expression := ExpressionBuilder().
		StartOfLine().
		CaptureWithName("foo").
		Then("a").
		Dot().
		EndCapture().
		Build()
	println(expression.String())
	println(expression.regexp.MatchString("ababababababa"))
	str := `aba`
	t.Log(expression.GetText(str, "foo"))
}

func TestModifier(t *testing.T) {
	r := regexp.MustCompile(`(?ism)CaSe`)
	res := r.MatchString("case")
	t.Log(res)

	expression := ExpressionBuilder().
		StartOfLine().
		Then("CaSe").
		Build()
	t.Log(expression.String())
	t.Log(expression.Test("CaSe"))
	t.Log(expression.Test("case"))
	t.Log(expression.Test("CASE"))
}

func TestGroup(t *testing.T) {
	str := `abcaAabbbxxxx`
	expression := ExpressionBuilder().
		StartOfLine().
		Anything().
		CaptWithName("foo").
		Then("a").Count(2).
		EndCapt().
		Anything().
		EndOfLine().
		WithAnyCase().
		Build()

	t.Log(expression.String())
	t.Log(expression.Test(str))
	t.Log(expression.GetText(str, "foo"))
	t.Log(expression.GetTextGroups(str, 0))
}

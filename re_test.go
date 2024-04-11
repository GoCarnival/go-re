/* Package re
 * @Author 砚池/Ivan
 * @Date 2024/04/10
 * @Description:
 */

package re

import "testing"

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
	str := `abacad`
	t.Log(expression.GetText(str, "foo"))
}

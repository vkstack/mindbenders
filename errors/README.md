# Errors
## Problem
Builtin errors is very limited in terms of wrapping. The basic wrapping which is offered by builtin loses content we wish to retain for long.

From the repository layer one passes an error of **mysql.MySQLError** type. which is wrapped at multipe layers of application[repository → usecase1→ usecase2 → super-usecase]. And in the super usecase one want's to do something if the cause is of **mysql.MySQLError** type.

## Solutions
wrapping can be done in the manner.
***e3 → e2 → e1 → e0***
e3 is the top level error and e0 is ground level error.
That is how you can do the wrapping.

```
import "gitlab.com/dotpe/mindbenders/errors"
.
,

e0 := errors.New("base-error")
e1 := errors.WrapMessage(e0, "wrapped e0")
e2 := errors.WrapMessage(e1, "wrapped e1")
e3 := errors.WrapMessage(e2, "wrapped e2")

errors.UnWrap(e3) // -> e2
errors.UnWrap(e2) // -> e1
errors.UnWrap(e0) // -> nil because e0 is the very base error

errors.Cause(e3) // -> e0
errors.Cause(e2) // -> e0
errors.Cause(e1) // -> e0
errors.Cause(e0) // -> e0
```
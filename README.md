# Goption

![goption](goption.png "Goption")

A way to handle nil values.

## Requirements

Use Go 1.18+ is a must. Goption uses Generics.

## Description

This is a copy of [Java Optional](https://docs.oracle.com/javase/8/docs/api/java/util/Optional.html). It contains almost all functions

- Empty
- Of
- Optional.Get
- Optional.IsPresent
- Optional.OrElseError
- Optional.OrElse

And others made specifically to go:

- Optional.MustGet

## SQL

You can use goption.Optional to sql package. Optional implements sql.Scanner and sql.Valuer

## JSON Marshall

Unfortunately `,omitempty` is not supported by now due how json.Marshal works. For more details can see [this link](https://github.com/golang/go/issues/11939).

_edit_: Now in Go version 1.24 [IsZero()](https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson) function is called by encoding/json when `,omitzero` tag is set.

## Boolean values

There is a special use for boolean type. By default, true and false will be valid when are unmarshal from json or text. The only option to handle as invalid bool (nil) is assigning it as a pointer:

```golang
type User struct {
	IsValid goption.Optional[*bool] // if false will be invalid
}

type Client struct {
	IsValid goption.Optional[bool] // if false will be valid
}
```

These values are taken as true:

- true
- 1
- on
- yes

And these other as false:

- false
- 0
- off
- no

This is the design decision I had to make. If any please create an issue if you have a better idea for this :D

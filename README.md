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

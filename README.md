# Goption

![goption](goption.png "Goption")

A way to handle nil values.

## Description

This is a copy of [Java Optional](https://docs.oracle.com/javase/8/docs/api/java/util/Optional.html). It contains almost all functions

- Empty
- Of
- Goptional.Get
- Goptional.IsPresent
- Goptional.OrElseError
- Goptional.OrElse

__Consider to use Go 1.18+. Goption uses Generics.__

## SQL

You can use goption.Optional to sql package. Optional implements sql.Scanner and sql.Valuer

## JSON Marshall

Unfortunately `,omitempty` is not supported by now due how json.Marshal works. For more details can see [this link](https://github.com/golang/go/issues/11939).

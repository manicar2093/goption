# Goption

![goption](goption.png "Goption")

__On development__

A way to handle nil

## Description

This is a copy of [Java Optional](https://docs.oracle.com/javase/8/docs/api/java/util/Optional.html). It contains almost all functions

- Empty
- Of
- Goptional.Get
- Goptional.IsPresent
- Goptional.OrElseError
- Goptional.OrElse

## JSON Marshall

Unfortunately `,omitempty` is not supported by now due how json.Marshal works. For more details can see [this link](https://github.com/golang/go/issues/11939).

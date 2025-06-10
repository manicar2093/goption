## v0.11.0 (2025-06-10)

### Feat

- add bool pointer handler for text unmarshal
- add bool pointer handler for json unmarshal
- add unmarshal false bool will be always valid data
- add string bool unmarshal in text and json

## v0.10.0 (2025-06-03)

### Feat

- add is zero implementation for go 1.24

## v0.9.3 (2025-01-09)

### Fix

- add slice support

## v0.9.2 (2024-12-10)

### Fix

- change regrex to identify uuid

## v0.9.1 (2024-11-06)

### Fix

- add sql valuer implementation
- add call of scan method if implemented

## v0.9.0 (2024-10-18)

### Feat

- add new way to identify strings and numbers for unmarshalText
- add unmarshalText with tests

### Fix

- change unmarshalText implementation

## v0.8.5 (2024-07-19)

### Fix

- add handling for numbers and floats

## v0.8.4 (2024-07-18)

### Fix

- add handling multiline's strings \r carriage return

## v0.8.3 (2024-07-17)

### Fix

- add handling multiline's strings to be unmarshal

## v0.8.2 (2024-06-13)

### Fix

- add support for custom types

## v0.8.1 (2024-06-12)

### Fix

- add native null to be handle in a proper way

## v0.8.0 (2024-05-24)

### Feat

- add util func to get bool from isValidData func
- change way to handle null on unmarshall json

### Fix

- add recognition of empty string in unmarshall json

## v0.7.0 (2024-05-24)

### Feat

- add new way to assign values on sql scan

## v0.6.0 (2024-04-17)

### Feat

- add slices cases detecting if they are empty

## v0.5.0 (2024-04-12)

### Feat

- update go version to 1.21

## v0.4.0 (2024-04-12)

### Feat

- add echo binder implementation

### Fix

- add echo binder
- add echo binder

## v0.3.1 (2023-05-02)

### Perf

- modify golang version (#6)

## v0.3.0 (2023-01-02)

### Feat

- add mustGet func

## v0.2.0 (2022-12-29)

### Feat

- add json and sql implementations (#3)

## v0.1.0 (2022-12-28)

### Feat

- add cz config file (#2)
- add goreleaser step (#1)
- add basic methods to optional

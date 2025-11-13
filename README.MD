# Ecto
Library for deserialization and validation of input data**

## Basic definitions

*Schema* — set of rules for transforming and validating input
data. It consists of tests. A schema of type *Atomic* also contains a converter function and can
provide a default value.

*Converter function* — function that transforms input data from
one type to another for a subsequent test application.

*Test* - part of the schema consisting of a predicate function and
an error template returned when the predicate fails.
Tests may apply to all data types or only a restricted subset.

*Predicate function* — function describing a property or relation
of a subject, which may be either true or false. Returns `true` /
`false`.

## Schema types

### Atomic
Intended for describing scalar Go types or types that do not require
recursive/internal processing (integers, booleans, strings, etc.). \
Error format: 

```json
["error1", "error2"]
```

### Struct
A composite schema that is represented as an associative list (dictionary),
where each key is a struct field and the value is the schema for that
field.\
Example:

``` text
1 Struct-schema 1 {
2     Email: (Atom-schema 1)
3     Password: (Atom-schema 2)
4     Meta: (Struct-schema 2)
5 }
```
Error format:
``` json
{"Email": ["error1", "error2"], "Meta": {"meta1": ["error"]}}
```

### List
A composite schema that applies a selected subschema to each element of
an array/list.\
Error format:

``` json
{"0": ["error"], "2": ["error"]}
```

### Optional Schema
A composite wrapper over any schema. The inner schema is applied only when data is present (not null).
Requires a pointer type.\
An error format is identical to *Atomic*.


The principle of composite schemas boils down to calling nested sub-schemas with additional options, 
while *Atomic* is a standalone schema and its mechanism is as follows:

<!--![Diagram Image Link](./atomic_schema.puml)-->

## Tests

- Common for all types:
  - **Required**: value must not be zero-value. If this test fails, all subsequent tests are skipped.\
    Error: `required`
  - **OneOf**: value must be one of the provided list.\
    Error: `must be one of %v`

- For integers
  - **Eq**: value must equal the specified integer\
    Error: `must be equal to %d`
  - **Min**: minimum value (inclusive)\
    Error: `must be %d minimum`
  - **Max**: maximum value (inclusive)\
    Error: `must be %d maximum`

- For floating-point numbers
  - **Min**: minimum value (inclusive)\
    Error: `must be %f minimum`
  - **Max**: maximum value (inclusive)\
    Error: `must be %f maximum`
  - **MaxPrecision**: limit on the number of digits after the decimal point\
    Error: `has more than %d precision digits`

- For strings
  - **Min**: minimum number of characters\
    Error: `must be at least %d characters long`
  - **Max**: maximum number of characters\
    Error: `must be at most %d characters long`
  - **Regex**: regular expression validation\
    Error: `must match regex %s`
  - **URL**: Checks if the value is a valid URL\
    Error: `invalid URL`
  - **Currency**: ISO-4217 currency standard\
    Error: `invalid ISO-4217 currency`
  - **IP**: Checks IPv4/IPv6\
    Error: `invalid IP address`
  - **Lang**: ISO-639-1 language code\
    Error: `invalid ISO-639-1 language code`
  - **Country**: ISO-3166 country code\
    Error: `invalid ISO-3166 country code`

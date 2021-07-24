# json-go
A no-import json library written in go for educational purposes.

**Please do NOT use this code for production purposes.**

The goal of this project is to better understand the `json` format while practising [TDD](https://en.wikipedia.org/wiki/Test-driven_development) and general good code practices. In addition, to do so without relying on the standard library conveniences.

The `json` format <https://www.json.org/> is made up of the types:
* Keywords (`true`, `false`, `null`)
* Numbers
* Strings
* Arrays
* Objects
* Whitespace

A future goal will be to be able to parse this sample:
```json
{
  "types": ["Objects", "Arrays", "Strings", "Numbers", "Keywords", "Whitespace"],
  "examples": {
    "object": {"key1": "value1", "key2": "value2"},
    "array": ["value1", "value2"],
    "string": "value",
    "number": 1000,
    "keyword": true
  }
}
```

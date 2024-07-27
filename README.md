# SLParser
An LARL parser generator (Work in Progress) -- DONE JUST FOR FUN; Not meant to be official or anything.


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Introduction](#introduction)


## Introduction

In order to use the SLParser, the desired grammar must be specified in two separate files with two distinct grammars and formats: one for the lexer and the other for the parser.



## Lexer Grammar

```ebnf
equal = "=" .
dot = "." .
newline = [ "\r" ] "\n" .
quote = "\"" .
backslash = "\\" .
pipe = "|" .
tab = "\t" .
ws = " " .
right_arrow = "->" .
skip = "skip" .
range = ".." .

id = "a".."z" { "a".."z" } .

// \u0000..\u0021, "\"", \u0023..\u005B, "\\", \u005D..\uFFFF
char
   = \u0000..\u0021 | \u0023..\u005B | \u005D..\uFFFF
   | backslash ( quote | backslash )
   .
```


```ebnf
Source = Rule { newline { newline } Rule }  EOF .

Rule = id ws equal ws Rhs ws dot [ SkipCls ] .

Rule = id newline tab equal ws Rhs { newline tab pipe ws Rhs } newline tab dot [ SkipCls ] .

Rhs
   = quote char quote
   | Range
   .

Range
   = quote char quote range quote char quote
   .

SkipCls
   = right_arrow skip
   .
```


`
equal = "=" .
dot = "." .
pipe = "|" .
newline = [ "\r" ] "\n" .
tab = "\t" .
ws = " " . -> skip
op_paren = "(" .
cl_paren = ")" .

uppercase_id = "A".."Z" { "a".."z" } .
lowercase_id = "a".."z" { "a".."z" } .


## Parser Grammar



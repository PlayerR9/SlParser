# SLParser
An LARL parser generator (Work in Progress) -- DONE JUST FOR FUN; Not meant to be official or anything.


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Introduction](#introduction)


## Introduction

In order to use the SLParser, the desired grammar must be specified in two separate files with two distinct grammars and formats: one for the lexer and the other for the parser.


# EbnfParser
A parser that parses a modified version of the EBNF grammar.


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Introduction](#introduction)


## Introduction



## Language Specification

### Lexical Elements

Here's a description of the various lexical elements that appears in the grammar.


***Identifiers***

An **identifier** is defined as a sequence of one or more letters of the English alphabet (i.e., from 'a' to 'z') that, optionally, can end with a sequence of one or more decimal digits (i.e., from '0' to '9').

This grammar distinguishes between two kinds of identifiers: uppercase and lowercase. A **lowercase identifier** is a type of identifies whose letters can only be lowercase letters while an **uppercase identifier** is a type of identifier whose letters can either be uppercase or lowercase letters. Finally, lowercase identifiers can use a single underscore character (`_`) to separate words.

For example, `foo` and `f` are valid lowercase identifiers while `Foo`, `FooBar`, and `F` are valid uppercase identifiers. On the contrary, `fo o`, `1bar`, `foo__bar`, `fooBar`, and so on are not valid identifiers.



***Symbols***

A **symbol** is a special character that appears in the grammar.


| **Punctuation** | **Name** | **Description** |
| :---: | :---: | :--- |
| `.` | dot | specifies the end of a rule. |


| **Brackets** | **Name** | **Description** |
| :---: | :---: | :--- |
| `(`, `)` | parentheses | specifies the start and end of a sub OR rule. |


| **Operators** | **Name** | **Description** |
| :---: | :---: | :--- |
| `=` | equal | separates the lhs from the rhs. |
| `\|` | pipe | exclusive or. |


| **Whitespace** | **Name** | **Description** |
| :---: | :---: | :--- |
| `\r\n`, `\n` | newline | separates multiple rules and/or lines. |
| `\t` | tab | indentation. |
| ` ` | ws | separates elements from each other. |

Spaces and tabs are ignored in the grammar and so, stuff like `a b` and `a  b` are equivalent.


### Source

***Overview***

In this context, the term "source" refers to the file containing the EBNF grammar.


***Syntax***

Here's the syntax of the source file:
```ebnf
Source = Rule { "\n" Rule } EOF .
```
Where:
- `Rule` refer to the rules of the grammar.
- `EOF` is a special symbol that indicates the end of the file. Thus, outside of the rules, nothing else is allowed.


In essence, a source file is a sequence of one or more rules (each of which is separated by one or more newline characters (`\n`)) that are read until the end of the file.


### Rule

***Overview***

A **rule** is the core of any grammar and it is used to describe how the grammar should be parsed.


***Syntax***

Here's the syntax of a rule:
```ebnf
Rule     = SlRule | MlRule .
SlRule   = uppercase_id "=" RhsCls "." .
MlRule   = uppercase_id "\n" LineRule "\n." .
LineRule = "=" RhsCls { "\n| "RhsCls } .
RhsCls   = Rhs { Rhs } .
```
Where:
- `uppercase_id` refers to an uppercase identifier.
- `Rhs` refers to the right-hand side of the rule.


In essence, a rule can either be a single-line rule or a multi-line rule. If it is a single-line rule, then the uppercase identifier is followed by an equal sign (`=`) and the right-hand side clause followed by the dot (`.`). On the other hand, if it is a multi-line rule, then the uppercase identifier is followed by a sequence of one or more right-hand side clauses preceded by a pipe (`|`). Each line is indented one level and the first one is the only one that stats with an equal sign (`=`) rather than a pipe. Finally, the dot (`.`) is written in a newline and indented one level as well.


***Examples***

Here are some examples of valid rules:

```ebnf
Color
   = red
   | green
   | blue
   .
```
This rule states that a color can either be "red", "green", or "blue".


```ebnf
Person = name age .
```
This rule states that a person has a name followed by an age.


### Right-hand Side

***Overview***

A **right-hand side** is the unit of the grammar and it specifies the individual atoms/units that make up a rule.


***Syntax***

Here's the syntax of a right-hand side:
```ebnf
Rhs        = Identifier | OrGroup .
Identifier = uppercase_id | lowercase_id .
OrGroup    = "(" OrExpr ")" .
OrExpr     = Identifier "|" Identifier { "|" Identifier } .
```

In essence, a right-hand side can either be an identifier or an OR group. An identifier is any lowercase or uppercase word while, an OR group is an OR expression that is surrounded by parentheses (`(` and `)`). Finally, an OR expression is a sequence of two or more identifiers separated by a pipe (`|`).



### Parsing




***Full Grammar***

```ebnf
equal = "=" .
dot = "." .
pipe = "|" .
newline = [ "\r" ] "\n" { [ "\r" ] "\n" } .
ws = " " | "\t" . -> skip
op_paren = "(" .
cl_paren = ")" .

uppercase_id = uppercase_word { uppercase_word } { digit } .
lowercase_id = lowercase_word { digit } .

fragment lowercase_word = "a".."z" { "a".."z" } . 
fragment uppercase_word = "A".."Z" { "a".."z" } .
fragment digit = "0".."9" .

Source = Source1 EOF .
Source1 = Rule .
Source1 = Rule newline Source1 .

Rule = uppercase_id equal RhsCls dot .
Rule = uppercase_id newline equal RhsCls RuleLine .
RuleLine = newline pipe RhsCls RuleLine .
RuleLine = newline dot  .

RhsCls = Rhs .
RhsCls = Rhs RhsCls .

Rhs = Identifier .
Rhs = op_paren OrExpr cl_paren .

OrExpr = Identifier pipe Identifier .
OrExpr = Identifier pipe OrExpr .

Identifier = uppercase_id .
Identifier = lowercase_id .
```




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



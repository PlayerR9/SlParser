lexer grammar SlLexer;

// Identifiers

UPPERCASE_ID : [A-Z]+([_][A-Z]+)* ;
LOWERCASE_ID : [a-z]+([A-Z][a-z]*)* ;


// Operators

MOD : [*?+];
PIPE : '|';


// SYMBOLS

// Punctuation

COLON : ':' ;
SEMICOLON : ';' ;


// Brackets

OP_PAREN : '(' ;
CL_PAREN : ')' ;


// Whitespace

WS : [ \t\r\n]+ -> skip ;
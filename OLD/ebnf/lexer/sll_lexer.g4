lexer grammar SlLexer;

// Identifiers

UPPERCASE_ID : [A-Z]+([_][A-Z]+)* ;


// Operators

MOD : [*?+];
PIPE : '|';


// SYMBOLS

// Punctuation

COLON : ':' ;
SEMICOLON : ';' ;
QUOTE : '\'' ;
BACK_SLASH : '\\' ;


// Brackets

OP_PAREN : '(' ;
CL_PAREN : ')' ;
OP_SQUARE : '[' ;
CL_SQUARE : ']' ;


// Whitespace

WS : [ \t\r\n]+ -> skip ;
lexer grammar Lexer;

// KEYWORDS

// Identifiers

LOWERCASE_ID : [a-z]+([A-Z][a-z]*)* ;
UPPERCASE_ID : [A-Z]+([_][A-Z]+)* ;


// SYMBOLS

// Punctuations

COLON : ':' ;
SEMICOLON : ';' ;

// Whitespaces

NEWLINE : ('\r'? '\n')+ ;
WS : [ \t]+ -> skip ;
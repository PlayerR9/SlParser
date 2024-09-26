lexer grammar Lexer;

// KEYWORDS

// Identifiers

LOWERCASE_ID : LOWERCASE ANY* ;
UPPERCASE_ID : UPPERCASE+ (UNDERSCORE NONLOWER+)* ;


// SYMBOLS

// Punctuations

COLON : ':' ;
SEMICOLON : ';' ;

// Whitespaces

NEWLINE : ('\r'? '\n')+ ;
WS : [ \t]+ -> skip ;


// FRAGMENTS

fragment ANY : [a-zA-Z0-9] ;
fragment NONUPPER : [a-z0-9] ;

fragment UPPERCASE : [A-Z] ;
fragment LOWERCASE : [a-z] ;
fragment UNDERSCORE : '_' ;
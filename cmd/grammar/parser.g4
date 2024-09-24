parser grammar Parser;

options {
   tokenVocab = Lexer;
};

source : rule (NEWLINE rule)* EOF ;
rule : LOWERCASE_ID COLON rhs+ SEMICOLON ;
rhs
   : UPPERCASE_ID
   | LOWERCASE_ID
   ;
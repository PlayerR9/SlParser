parser grammar Parser;

options {
   tokenVocab = Lexer;
};

source : NEWLINE statement (NEWLINE statement)* EOF ;

statement
   : LIST_COMPREHENSION
   | PRINT_STMT
   ;
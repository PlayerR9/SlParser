parser grammar Parser;

options {
   tokenVocab = Lexer;
};

source : statement (NEWLINE statement)* EOF ;

statement
   : LIST_COMPREHENSION
   | PRINT_STMT
   ;
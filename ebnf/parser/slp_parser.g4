parser grammar SlParser;

options {
   tokenVocab = SlLexer;
};

source : rule+ EOF ;

rule : LOWERCASE_ID COLON expression? SEMICOLON ;

expression : andExpr (PIPE andExpr)* ;

andExpr : modExpr+ ;

modExpr : rhs MOD? ;

rhs : UPPERCASE_ID | LOWERCASE_ID | OP_PAREN expression CL_PAREN ;
parser grammar SlParser;

options {
   tokenVocab = SlLexer;
};

source : rule+ EOF ;

rule : UPPERCASE_ID COLON expression? SEMICOLON ;

literal : QUOTE (any char | BACK_SLASH (BACK_SLASH | QUOTE)) QUOTE ;

expression : andExpr (PIPE andExpr)* ;

andExpr : modExpr+ ;

modExpr : rhs MOD? ;

rhs : UPPERCASE_ID | OP_PAREN expression CL_PAREN | OP_SQUARE ()* CL_SQUARE | literal ;
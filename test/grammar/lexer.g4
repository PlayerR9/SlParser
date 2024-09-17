lexer grammar Lexer;

LIST_COMPREHENSION : 'sq = [x * x for x in range(10)]' ;
PRINT_STMT : 'sq' ;

NEWLINE : ('\r'? '\n')+ ;

COMMENT : '#' .*? '\n' -> skip ;
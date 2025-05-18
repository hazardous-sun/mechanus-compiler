# Mechanus Derivation Tree

```
<G> ::= '{' <BODY> '}' <ID> 'Construct'

<BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect'
<BODY> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect'
<BODY> ::= <BODY_REST>

<BODY_REST> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' <BODY_REST>
<BODY_REST> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect' <BODY_REST>
<BODY_REST> ::= ε

<TYPE> ::= 'Nil'
<TYPE> ::= 'Gear'
<TYPE> ::= 'Tensor'
<TYPE> ::= 'State'
<TYPE> ::= 'Monodrone'
<TYPE> ::= 'Omnidrone'

<CMDS> ::= <CMD> <CMDS_REST>

<CMDS_REST> ::= '\n' <CMDS>
<CMDS_REST> ::= ε

<CMD> ::= <CMD_IF>
<CMD> ::= <CMD_FOR>
<CMD> ::= <CMD_DECLARATION>
<CMD> ::= <CMD_ASSIGNMENT>
<CMD> ::= <CMD_RECEIVE>
<CMD> ::= <CMD_SEND>

<CMD_IF> ::= '{' <CMDS> '}' 'if' <CONDITION>
<CMD_IF> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'if' <CONDITION>  
<CMD_IF> ::= <CMD_ELIF> '{' <CMDS> '}' 'if' <CONDITION>

<CMD_ELIF> ::= '{' <CMDS> '}' 'elif' <CONDITION>
<CMD_ELIF> ::= <CMD_ELIF_REST>
<CMD_ELIF_REST> ::= '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST>
<CMD_ELIF_REST> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST>
<CMD_ELIF_REST> ::= ε

<CMD_FOR> ::= '{' <CMDS> '}' 'for' <CONDITION>

<CMD_DECLARATION> ::= <E> '=:' <VAR>

<CMD_ASSIGNMENT> ::= <E> '=' <VAR> 

<CMD_RECEIVE> ::= '(' <VAR> ')' 'Receive'

<CMD_SEND> ::= '(' <E> ')' 'Send'

<CONDITION> ::= <E> '>' <E> 
<CONDITION> ::= <E> '>=' <E> 
<CONDITION> ::= <E> '<>' <E> 
<CONDITION> ::= <E> '<=' <E> 
<CONDITION> ::= <E> '<' <E> 
<CONDITION> ::= <E> '==' <E>

<E> ::= <T> <E_REST>

<E_REST> ::= '+' <T> <E_REST>
<E_REST> ::= '-' <T> <E_REST>
<E_REST> ::= ε

<T> ::= <F> <T_REST>

<T_REST> ::= '*' <F> <T_REST>
<T_REST> ::= '/' <F> <T_REST>
<T_REST> ::= '%' <F> <T_REST>
<T_REST> ::= ε

<F> ::= -<F> | <X>

<X> ::= '(' <E> ')'
<X> ::= [0-9]+('.'[0-9]+)
<X> ::= <VAR>

<PARAMETERS> ::= <ID> ':' <TYPE>
<PARAMETERS> ::= <ID> ':' <TYPE> <EXTRA_PARAMETERS>

<EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE>
<EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> ',' <PARAMETERS>

<ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+

<TEXT_WITH_NUMBERS> ::= (([A-Z]|[a-z])*(_|[0-9])*)+
<TEXT_WITHOUT_NUMBERS> ::= (([A-Z]|[a-z])+(_)*)+
```
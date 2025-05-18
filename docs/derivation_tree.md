# Mechanus Derivation Tree

```
<G> ::= '{' <BODY> '}' <TEXT_WITHOUT_NUMBERS> 'Construct'

<BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
<BODY> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
<BODY> ::= <BODY_REST>

<BODY_REST> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect' <BODY_REST>
<BODY_REST> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect' <BODY_REST>
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

<CMD_DECLARATION> ::= <VAR> ':=' <E>

<CMD_ASSIGNMENT> ::= <VAR> '=' <E>

<CMD_RECEIVE> ::= 'Receive' '(' <VAR> ')'

<CMD_SEND> ::= 'Send' '(' <E> ')'

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

<PARAMETERS> ::= <TEXT_WITH_NUMBERS> ':' <TYPE>
<PARAMETERS> ::= <TEXT_WITH_NUMBERS> ':' <TYPE> <EXTRA_PARAMETERS>

<EXTRA_PARAMETERS> ::= ',' <TEXT_WITH_NUMBERS> ':' <TYPE>
<EXTRA_PARAMETERS> ::= ',' <TEXT_WITH_NUMBERS> ':' <TYPE> ',' <PARAMETERS>

<TEXT_WITH_NUMBERS> ::= (([A-Z]|[a-z])+(_|[0-9])*)+
<TEXT_WITHOUT_NUMBERS> ::= (([A-Z]|[a-z])+(_)*)+
```
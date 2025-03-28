# Mechanus Derivation Tree

```
<G> ::= '{' <BODY> '}' <TEXT_WITHOUT_NUMBERS> 'Construct'

<BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
<BODY> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
<BODY> ::= <BODY> '{' <CMDS> '}' '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
<BODY> ::= <BODY> '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'

<TYPE> ::= 'Nil'
<TYPE> ::= 'Gear'
<TYPE> ::= 'Tensor'
<TYPE> ::= 'State'
<TYPE> ::= 'Monodrone'
<TYPE> ::= 'Omnidrone'

<CMDS> ::= <CMD>
<CMDS> ::= <CMD> '\n' <CMDS>

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
<CMD_ELIF> ::= <CMD_ELIF> '{' <CMDS> '}' 'elif' <CONDITION>
<CMD_ELIF> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'elif' <CONDITION>

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

<E> ::= <E> + <T>
<E> ::= <E> - <T>
<E> ::= <T>

<T> ::= <T> * <F>
<T> ::= <T> / <F>
<T> ::= <T> % <F>
<T> ::= <F>

<F> ::= -<F>
<F> ::= <X>

<X> ::= '(' <E> ')'
<X> ::= [0-9]+('.'[0-9]+)
<X> ::= <VAR>

<PARAMETERS> ::= <TEXT_WITH_NUMBERS> ':' <TYPE>
<PARAMETERS> ::= <TEXT_WITH_NUMBERS> ':' <TYPE> <EXTRA_PARAMETERS>

<EXTRA_PARAMETERS> := ',' <TEXT_WITH_NUMBERS> ':' <TYPE>
<EXTRA_PARAMETERS> := ',' <TEXT_WITH_NUMBERS> ':' <TYPE> ',' <PARAMETERS>

<TEXT_WITH_NUMBERS> ::= (([A-Z]|[a-z])+(_|[0-9])*)+
<TEXT_WITHOUT_NUMBERS> ::= (([A-Z]|[a-z])+(_)*)+
```
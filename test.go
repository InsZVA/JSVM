package main

import (
    "./lex"
    "./syntax"
)

func main() {
    peeker := lex.NewPeeker([]byte("a =1+2+3/6"))
    var s syntax.TokenStream
    for t, _ := peeker.ReadToken();t != nil;t, _ = peeker.ReadToken() {
        s.Tokens = append(s.Tokens, t)
    }
    tree, err := syntax.UnEndableParserList[syntax.ASSIGN](&s)
    if err != nil {
        panic(err)
    }
    tree.Print()
}
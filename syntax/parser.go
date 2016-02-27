package syntax

import (
    "../lex"
	"errors"
    "container/list"
    "fmt"
)

const (
    ASSIGN    =         iota
    EXP
)

type TokenStream struct {
    Tokens              []*lex.Token
    Pos                 int
}

func (s *TokenStream) ReadToken() *lex.Token {
    if s.Pos >= len(s.Tokens) {
        return nil
    }
    s.Pos++
    return s.Tokens[s.Pos - 1]
}

func (s *TokenStream) Fork() *TokenStream {
    return &TokenStream {
        Tokens:         s.Tokens,
        Pos:            s.Pos,
    }
}

type SyntaxTreeNode struct {
    LNode, RNode        *SyntaxTreeNode
    Name                string
    PrintName           string
}

func (s *SyntaxTreeNode) Print() {
    if s == nil {
        return
    }
    s.PrintName = s.Name
    queue := list.New()
    queue.PushBack(s)
    for queue.Len() != 0 {
        front := queue.Front()
        fmt.Println(front.Value.(*SyntaxTreeNode).PrintName, " ")
        if front.Value.(*SyntaxTreeNode).LNode != nil {
            front.Value.(*SyntaxTreeNode).LNode.PrintName = front.Value.(*SyntaxTreeNode).PrintName + "\\" + front.Value.(*SyntaxTreeNode).LNode.Name
            queue.PushBack(front.Value.(*SyntaxTreeNode).LNode)
        }
        if front.Value.(*SyntaxTreeNode).RNode != nil {
            front.Value.(*SyntaxTreeNode).RNode.PrintName = front.Value.(*SyntaxTreeNode).PrintName + "\\" + front.Value.(*SyntaxTreeNode).RNode.Name
            queue.PushBack(front.Value.(*SyntaxTreeNode).RNode)
        }
        queue.Remove(front)
    }
}

type UnEndableParser func(*TokenStream) (*SyntaxTreeNode, error)

var UnEndableParserList map[int]UnEndableParser

func EXPParser(s *TokenStream) (*SyntaxTreeNode, error) {
    fork := s.Fork()
    tk1 := fork.ReadToken()
    if tk1 == nil {
        return NullProduce()(s)
    }
    tk2 := fork.ReadToken()
    if tk2 == nil {
        return TokenProduce(tk1.Type)(s)
    } else if tk2.Type == lex.COMPUTE{
        return SumProduce(TokenProduce(tk1.Type), TokenProduce(lex.COMPUTE), EXPParser)(s)
    }
    return TokenProduce(tk1.Type)(s)
}

func init() {
    UnEndableParserList = make(map[int]UnEndableParser)
    /*
    UnEndableParserList[EXP] = NullProduce()
    UnEndableParserList[EXP] = OrProduce(SumProduce(UnEndableParserList[EXP], TokenProduce(lex.COMPUTE), UnEndableParserList[EXP]), TokenProduce(lex.IDENTIFIER))*/
    UnEndableParserList[EXP] = EXPParser
    UnEndableParserList[ASSIGN] = SumProduce(TokenProduce(lex.IDENTIFIER), TokenProduce(lex.EQUAL), UnEndableParserList[EXP])
}

func NullProduce() (UnEndableParser) {
    return func(s *TokenStream) (*SyntaxTreeNode, error) {
        return &SyntaxTreeNode {
            Name:       "NullNode",
        }, nil
    }
}

func TokenProduce(t int) (UnEndableParser) {
    return func(s *TokenStream) (*SyntaxTreeNode, error) {
        tk := s.ReadToken()
        if tk.Type == t {
            return &SyntaxTreeNode {
                Name:       tk.Chars + "Node",
            }, nil
        } else {
            return nil, errors.New("Unexpected Token:" + tk.Chars)
        }
    }
}

func OrProduce(p1 UnEndableParser, p2 UnEndableParser) (UnEndableParser) {
    return func(s *TokenStream) (*SyntaxTreeNode, error) {
        var err error
        fork := s.Fork()
        if r, err := p1(fork);err == nil {
            s.Pos = fork.Pos
            return r, err
        }
        fork = s.Fork()
        if r, err := p2(fork);err == nil {
            s.Pos = fork.Pos
            return r, err
        } 
        return nil, err
    }
}

func AndProduce(p1 UnEndableParser, p2 UnEndableParser) (UnEndableParser) {
    return func(s *TokenStream) (*SyntaxTreeNode, error) {
        var err error
        var l, r *SyntaxTreeNode
        fork := s.Fork()
        if l, err = p1(fork);err != nil {
            return nil, err
        }
        if r, err = p2(fork);err != nil {
            return nil, err
        } 
        s.Pos = fork.Pos
        return &SyntaxTreeNode {
            Name:       "AndNode",
            LNode:      l,
            RNode:      r,
        }, nil
    }
}

func SumProduce(pArray ...UnEndableParser) (UnEndableParser) {
    pResult := NullProduce()
    for _, p := range pArray {
        pResult = AndProduce(pResult, p)
    }
    return pResult
}
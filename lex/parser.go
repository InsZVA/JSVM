package lex

import (
    "errors"
)

const (
    IDENTIFIER =   iota
    DISCRIPTION
    COMPUTE
    FOR
    IF
    ELSE
    BRACKET
    EQUAL
)

type Token struct {
    Type        int
    Chars       string
}

type Peeker struct {
    Raw         []byte
    RPos        int
    Line        int
}

func (p *Peeker) ReadToken() (*Token, error) {
    var buff string
    start:
    if p.RPos >= len(p.Raw) {
        if buff == "" {
            return nil, errors.New("EOF")
        } else {
            goto end
        }
    }
    switch p.Raw[p.RPos] {
        case ' ', '\t':
        p.RPos++
        if buff == "" {
            goto start
        } else {
            goto end
        }
        case '\n':
        p.RPos++
        p.Line++
        if buff == "" {
            goto start
        } else {
            goto end
        }
        case '{', '}', '+', '-', '*', '/', '^', '|', '&', '=':
        if buff == "" {
            buff = string(p.Raw[p.RPos])
            p.RPos++
            goto end
        } else {
            goto end
        }
        default:
        buff += string(p.Raw[p.RPos])
        p.RPos++
        goto start
    }
    end:
    var _type int
    switch buff {
        case "var":
        _type = DISCRIPTION
        case "for":
        _type = FOR
        case "if":
        _type = IF
        case "else":
        _type = ELSE
        case "+", "-", "*", "/", "&", "|", "^":
        _type = COMPUTE
        case "{", "}":
        _type = BRACKET
        case "=":
        _type = EQUAL
        case "":
        return nil, errors.New("Parse finished")
        default:
        _type = IDENTIFIER
    }
    return &Token {
        Type:       _type,
        Chars:      buff,
    }, nil
}

func NewPeeker(bytes []byte) *Peeker {
    new := make([]byte, len(bytes))
    copy(new, bytes)
    return &Peeker {
        Line:       1,
        Raw:        new,
    }
}
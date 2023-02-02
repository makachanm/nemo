package nmcore

import (
	"strings"
)

type NLexer struct {
	internalCounter BraketCounter
}

type NParser struct {
	internalCounter BraketCounter
	funcStack       BlockStack
}

func NewNLexer() NLexer {
	return NLexer{
		internalCounter: MakeBraketCounter(),
	}
}

func NewNParser() NParser {
	return NParser{
		internalCounter: MakeBraketCounter(),
		funcStack:       NewBlockStack(),
	}
}

func (l *NLexer) Toknize(input string, tokenmap map[string]TokMapElement) []Block {
	var totalLengths = uint64(len(input))
	var pointer uint64 = 0

	var processedBlocks []Block

	for pointer < totalLengths {
		currentChar := GetChar(&input, int(pointer))
		tokenVal, isToken := tokenmap[currentChar]

		if isToken && tokenVal.Token != TokenIgnore {
			AppendSingleBlock(&processedBlocks, GenerateBlock(tokenVal.Token, tokenVal.MatchChar))
			pointer++
		} else if tokenVal.Token == TokenIgnore {
			var startPos, endPos uint64
			startPos = pointer + 1
			currentPos := startPos
			isBreak := false

			for currentPos < totalLengths && !isBreak {
				if tokenmap[GetChar(&input, int(currentPos))].Token == TokenIgnore {
					isBreak = true
				}

				currentPos++
			}
			endPos = currentPos

			AppendSingleBlock(&processedBlocks, GenerateBlock(TokenString, GetCharfromRange(&input, int(startPos), int(endPos-1))))
			pointer = endPos

		} else {
			var startPos, endPos uint64
			currentPos := pointer
			isBreak := false

			startPos = pointer
			for currentPos < totalLengths && !isBreak {
				_, isBreak = tokenmap[GetChar(&input, int(currentPos))]
				currentPos++
			}
			endPos = currentPos - 1

			if totalLengths == (pointer + 1) { //if is last char of input?
				endPos++
			}

			AppendSingleBlock(&processedBlocks, GenerateBlock(TokenString, GetCharfromRange(&input, int(startPos), int(endPos))))
			pointer = endPos

		}
	}

	AppendSingleBlock(&processedBlocks, GenerateBlock(TokenEol, ""))
	return processedBlocks
}

func (p *NParser) Parse(input []Block) ExprNode {
	var originNode = MakeExprNode(TypeSection, nil)
	var pointer uint64 = 0

	for input[pointer].Token != TokenEol {
		p.funcStack.BlockPush(input[pointer])

		switch p.funcStack.Seek().Token {
		case TokenString:
			if p.funcStack.Length() <= 1 {
				stringBlock := p.funcStack.BlockPop()
				originNode.SingleInsert(MakeExprNode(TypeString, []string{stringBlock.Item}))
			}

		case TokenBlockStart:
			p.internalCounter.IncOpen()

		case TokenBlockEnd:
			p.internalCounter.IncClose()
			if p.internalCounter.IsBlocked() {
				rsnode := parseFuncStack(&p.funcStack)
				originNode.SingleInsert(rsnode)
			}
		}

		pointer++

	}

	return originNode
}

func parseFuncStack(input *BlockStack) ExprNode {
	var funcNode ExprNode

	if !(input.Length() > 1) && input.Seek().Token == TokenString { //check it is singe string token
		funcNode = MakeExprNode(TypeString, []string{input.BlockPop().Item})
		return funcNode
	}

	stackCtx := input.Clear()
	stackCtx = stackCtx[2:(len(stackCtx) - 1)] //remove unused symbols

	funcNode = parseFunc(stackCtx[0].Item)
	RemoveElementBlockArray(&stackCtx, 0)

	var innerFnStacks []BlockStack
	var symCounter = MakeBraketCounter()
	var innerfnsPointer uint64 = 0
	var fstart, fend uint64 = 0, 0

	//parse inner blocks to pieces
	for len(stackCtx) > 0 {
		switch stackCtx[innerfnsPointer].Token {
		case TokenFunc:
			if !symCounter.IsBlocked() {
				innerfnsPointer++
			} else {
				fstart = innerfnsPointer
				innerfnsPointer++
			}

		case TokenBlockStart:
			symCounter.IncOpen()
			innerfnsPointer++

		case TokenBlockEnd:
			symCounter.IncClose()
			if symCounter.IsBlocked() {
				fend = innerfnsPointer

				nstack := NewBlockStack()
				nstack.AppendArray(&stackCtx, int(fstart), int(fend))
				innerFnStacks = append(innerFnStacks, nstack)
				symCounter.Clear()
				innerfnsPointer = 0
			} else if len(stackCtx) <= 1 {
				RemoveElementBlockArray(&stackCtx, 0)
			} else {
				innerfnsPointer++
			}

		default:
			if symCounter.IsBlocked() || symCounter.IsZero() {
				fend = innerfnsPointer

				nstack := NewBlockStack()
				nstack.AppendArray(&stackCtx, int(fstart), int(fend))
				innerFnStacks = append(innerFnStacks, nstack)
				symCounter.Clear()
				innerfnsPointer = 0

			} else {
				innerfnsPointer++
			}

		}
	}

	if len(innerFnStacks) > 0 {
		for _, ptagetStack := range innerFnStacks {
			inpNode := parseFuncStack(&ptagetStack)
			funcNode.SingleInsert(inpNode)
		}
	}

	return funcNode
}

func parseFunc(input string) ExprNode {
	splitedD := strings.Split(input, " ")
	fnName := splitedD[0]
	fnCtx := []string{strings.Join(splitedD[1:], " ")}
	fnArgs := make(map[string]string)

	if strings.Contains(fnName, "(") && strings.Contains(fnName, ")") {
		splitPoint := strings.Index(fnName, "(") + 1
		argsctx := fnName[splitPoint:(len(fnName) - 1)]
		qname := fnName[:(splitPoint - 1)]

		ls := strings.Split(argsctx, ",")
		for _, ctxL := range ls {
			sptArg := strings.Split(ctxL, "=")
			fnArgs[sptArg[0]] = strings.Join(sptArg[1:], "=")
		}

		return MakeFunctionNode(qname, fnArgs, fnCtx)
	}

	return MakeFunctionNode(fnName, fnArgs, fnCtx)
}

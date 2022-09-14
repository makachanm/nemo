package nemomark

import (
	"strings"
)

type Lexer struct {
}

type Parser struct {
	funcBuffer BlockStack
	bcounter   BraketCounter
}

func NewLexer() Lexer {
	return Lexer{}
}

func NewParser() Parser {
	return Parser{
		funcBuffer: NewBlockStack(),
		bcounter:   MakeBraketCounter(),
	}
}

func returnString(s *string, idx int) string {
	return string((*s)[idx])
}

func addBlock(b *[]Block, t Token, i string) []Block {
	return append(*b, Block{token: t, item: i})
}

func appendEol(b *[]Block) []Block {
	return append(*b, Block{token: TokenEol})
}

func (l *Lexer) Tokenize(input string, tokenmap map[string]Token) []Block {
	var Blocks []Block
	var pointer = 0
	var lengths = len(input)

	for pointer < lengths {
		currentString := returnString(&input, pointer)
		tokenVal, keyExist := tokenmap[currentString]

		if keyExist {
			if tokenVal == TokenIgnore {
				iidx := pointer
				iisBreak := tokenVal == TokenIgnore

				for iisBreak {
					//Check idx is smaller than lengths. make break when string pointer is encountered eol.
					if (iidx + 1) == lengths {
						iidx++ //increase pointer
						break
					}
					t, r := tokenmap[returnString(&input, iidx+1)]
					iisBreak = !(r && t == TokenIgnore)
					iidx++
				}

				Blocks = addBlock(&Blocks, TokenString, input[(pointer+1):iidx])
				pointer = (iidx - pointer) + pointer + 1

			} else {
				//Tokenize current string.
				Blocks = addBlock(&Blocks, tokenVal, "")
				pointer++
			}

		} else {
			//Make current text pointer to normal string.
			idx := pointer //String tokenizer pointer
			isBreak := keyExist
			for !isBreak { //loop until enconter another token.
				//Check idx is smaller than lengths. make break when string pointer is encountered eol.
				if (idx + 1) == lengths {
					idx++ //increase pointer
					break
				}
				//Check token.
				_, isBreak = tokenmap[returnString(&input, idx+1)]
				idx++ //increase pointer (not eol)
			}
			//Cut string and tokenize. [Strat of Stirng to End of String]
			Blocks = addBlock(&Blocks, TokenString, input[pointer:idx])
			pointer = (idx - pointer) + pointer
		}
	}
	//Append EOL token.
	Blocks = appendEol(&Blocks)

	return Blocks
}

func (p *Parser) Parse(input *[]Block) ExprNode {
	var Blocks = *input
	var Head = MakeExprNode(TypeSection, nil) //Head of tree.
	var pointer = 0
	var lengths = len(Blocks)

	var Stack = NewBlockStack() //Stack for StackParse.

	for pointer < lengths {
		Stack.blockPush(Blocks[pointer])
		parsed := p.stackParse(&Stack)
		if parsed.Context != nil || parsed.FuncContext.FucntionName != "" {
			// Check Node's context is not a nil and Node is valid function node.
			Head.singleInsert(parsed)
		}
		pointer++
	}

	return Head
}

func (p *Parser) stackParse(input *BlockStack) ExprNode {
	var object ExprNode

	switch input.seek().token {
	case TokenString:
		if p.funcBuffer.length() > 1 { //Check another block is existed in stack.
			p.funcBuffer.blockPush(input.blockPop())
			break
		}
		item := input.blockPop().item

		stringTh := []string{item}
		object = MakeExprNode(TypeString, stringTh)

	case TokenBlockStart:
		openToken := input.blockPop()
		p.bcounter.incOpen()
		p.funcBuffer.blockPush(openToken)

	case TokenBlockEnd:
		closeToken := input.blockPop()
		p.bcounter.incClose()
		p.funcBuffer.blockPush(closeToken)

		funcParsed, isParsed := p.funcParse(&p.funcBuffer)

		if isParsed {
			object = funcParsed
		}

	default:
		p.funcBuffer.blockPush(input.blockPop())
		funcParsed, isParsed := p.funcParse(&p.funcBuffer)

		if isParsed {
			object = funcParsed
		}
	}

	return object
}

func (p *Parser) funcParse(input *BlockStack) (ExprNode, bool) {
	var object ExprNode

	frame := *(input)

	//Check function have minimal items.
	if !(frame.length() >= 4) {
		return object, false
	}

	//Check function is closed.
	if !(frame.seek().token == TokenBlockEnd) {
		return object, false
	}

	//Check open token and close token is all pushed to stack.
	if !(p.bcounter.isBlocked()) {
		return object, false
	}

	//Check function is starting with function symbol, Check function open token is existed in right place.
	//Clear buffer when token block position is not correct.
	if !(frame.blockList[0].token == TokenFunc && frame.blockList[1].token == TokenBlockStart) {
		input.clear()
		return object, false
	}

	//Is context string is Exist?
	if frame.blockList[2].token == TokenString {
		funcData := frame.blockList[2]               //Get context.
		splited := strings.Split(funcData.item, " ") //Split string to parse function name & args

		funcName := splited[0]                         //Get function name.
		pFuncContext := strings.Join(splited[1:], " ") //Join remain strings to one context value.

		var fnargs = make(map[string]string)

		if strings.ContainsAny(funcName, "(") && strings.ContainsAny(funcName, ")") {
			startpoint := strings.Index(funcName, "(")
			endpoint := strings.Index(funcName, ")")

			argstr := funcName[(startpoint + 1):endpoint]
			funcName = funcName[0:startpoint]

			argsSplited := strings.Split(argstr, ",")
			for _, value := range argsSplited {
				if strings.ContainsAny(value, "=") {
					args := strings.Split(value, "=")
					carg := args[1:]
					fnargs[args[0]] = strings.Join(carg, "=")
				}
			}
		}

		object = MakeFunctionNode(funcName, fnargs, []string{pFuncContext})

		//If inner-function is existed?
		if frame.blockList[3].token == TokenFunc {
			var innerFuncStack BlockStack
			innerFuncStack.appendStack(frame, 3, frame.length()-1)

			p.bcounter.decCounter()
			obj, innerft := p.parseInnerBlock(innerFuncStack)
			if innerft {
				object.insert(obj)
			}
		}

		input.clear()

	}

	p.bcounter.decCounter()

	return object, true
}

func (p *Parser) parseInnerBlock(input BlockStack) ([]ExprNode, bool) {
	var position = make([]int, 0)       //Position of deffunc
	var parseList []BlockStack          //Parsed Innerfuncs
	var bkcounter = MakeBraketCounter() //Counter for parse innerfunc braket
	var object []ExprNode

	for idx, item := range input.blockList {
		if item.token == TokenFunc {
			position = append(position, idx)
		}
	}

	//Return false when funcdef is not found.
	if len(position) == 0 || len(position) < 1 {
		return object, false
	}

	pointer := 0
	ipos := 0
	for pointer < len(input.blockList) {
		switch input.blockList[pointer].token {
		case TokenBlockStart:
			bkcounter.incOpen()
		case TokenBlockEnd:
			bkcounter.incClose()
		}

		if bkcounter.isBlocked() && (bkcounter.open > 0 && bkcounter.close > 0) {
			end := pointer + 1
			start := position[ipos]

			innerFuncStk := NewBlockStack()
			innerFuncStk.appendStack(input, start, end)

			parseList = append(parseList, innerFuncStk)
			bkcounter.clear()

			for pi, t := range position {
				if t > pointer {
					ipos = pi
					pointer = position[pi]
					innerFuncStk.clear()
					break
				}
			}

		}

		pointer++
	}

	for _, context := range parseList {
		pased, isPased := p.funcParse(&context)
		if isPased {
			object = append(object, pased)
		}
	}

	return object, true
}

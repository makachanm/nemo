package nemomark

import (
	nmcore "nemo/nemomark/core"
	"strings"
)

type Lexer struct {
}

type Parser struct {
	funcBuffer nmcore.BlockStack
	bcounter   nmcore.BraketCounter
}

func NewLexer() Lexer {
	return Lexer{}
}

func NewParser() Parser {
	return Parser{
		funcBuffer: nmcore.NewBlockStack(),
		bcounter:   nmcore.MakeBraketCounter(),
	}
}

func returnString(s *string, idx int) string {
	return string((*s)[idx])
}

func addBlock(b *[]nmcore.Block, t nmcore.Token, i string) []nmcore.Block {
	return append(*b, nmcore.Block{Token: t, Item: i})
}

func appendEol(b *[]nmcore.Block) []nmcore.Block {
	return append(*b, nmcore.Block{Token: nmcore.TokenEol})
}

func (l *Lexer) Tokenize(input string, tokenmap map[string]nmcore.Token) []nmcore.Block {
	var Blocks []nmcore.Block
	var pointer = 0
	var lengths = len(input)

	for pointer < lengths {
		currentString := returnString(&input, pointer)
		tokenVal, keyExist := tokenmap[currentString]

		if keyExist {
			if tokenVal == nmcore.TokenIgnore {
				iidx := pointer
				iisBreak := tokenVal == nmcore.TokenIgnore

				for iisBreak {
					//Check idx is smaller than lengths. make break when string pointer is encountered eol.
					if (iidx + 1) == lengths {
						iidx++ //increase pointer
						break
					}
					t, r := tokenmap[returnString(&input, iidx+1)]
					iisBreak = !(r && t == nmcore.TokenIgnore)
					iidx++
				}

				if len(Blocks) > 2 && Blocks[len(Blocks)-1].Token == nmcore.TokenString {
					previtem := Blocks[len(Blocks)-1].Item
					curitem := input[(pointer + 1):iidx]

					Blocks[len(Blocks)-1].Item = previtem + curitem
					pointer = (iidx - pointer) + pointer + 1
				} else {
					Blocks = addBlock(&Blocks, nmcore.TokenString, input[(pointer+1):iidx])
					pointer = (iidx - pointer) + pointer + 1
				}

			} else {
				//Tokenize current string.
				Blocks = addBlock(&Blocks, tokenVal, "")
				pointer++
			}

		} else {
			//Make current text pointer to normal string.
			idx := pointer //String tokenizer pointer
			isBreak := keyExist
			for !isBreak { //loop until enconter another Token.
				//Check idx is smaller than lengths. make break when string pointer is encountered eol.
				if (idx + 1) == lengths {
					idx++ //increase pointer
					break
				}
				//Check Token.
				_, isBreak = tokenmap[returnString(&input, idx+1)]
				idx++ //increase pointer (not eol)
			}
			//Cut string and tokenize. [Strat of Stirng to End of String]
			Blocks = addBlock(&Blocks, nmcore.TokenString, input[pointer:idx])
			pointer = (idx - pointer) + pointer
		}
	}
	//Append EOL Token.
	Blocks = appendEol(&Blocks)

	return Blocks
}

func (p *Parser) Parse(input *[]nmcore.Block) nmcore.ExprNode {
	var Blocks = *input
	var Head = nmcore.MakeExprNode(nmcore.TypeSection, nil) //Head of tree.
	var pointer = 0
	var lengths = len(Blocks)

	var Stack = nmcore.NewBlockStack() //Stack for StackParse.

	for pointer < lengths {
		Stack.BlockPush(Blocks[pointer])
		parsed := p.stackParse(&Stack)
		if parsed.Context != nil || parsed.FuncContext.FunctionName != "" {
			// Check Node's context is not a nil and Node is valid function node.
			Head.SingleInsert(parsed)
		}
		pointer++
	}

	return Head
}

func (p *Parser) stackParse(input *nmcore.BlockStack) nmcore.ExprNode {
	var object nmcore.ExprNode

	switch input.Seek().Token {
	case nmcore.TokenString:
		if p.funcBuffer.Length() > 1 { //Check another block is existed in stack.
			p.funcBuffer.BlockPush(input.BlockPop())
			break
		}
		item := input.BlockPop().Item

		stringTh := []string{item}
		object = nmcore.MakeExprNode(nmcore.TypeString, stringTh)

	case nmcore.TokenBlockStart:
		openToken := input.BlockPop()
		p.bcounter.IncOpen()
		p.funcBuffer.BlockPush(openToken)

	case nmcore.TokenBlockEnd:
		closeToken := input.BlockPop()
		p.bcounter.IncClose()
		p.funcBuffer.BlockPush(closeToken)

		funcParsed, isParsed := p.funcParse(&p.funcBuffer)

		if isParsed {
			object = funcParsed
		}

	default:
		p.funcBuffer.BlockPush(input.BlockPop())
		funcParsed, isParsed := p.funcParse(&p.funcBuffer)

		if isParsed {
			object = funcParsed
		}
	}

	return object
}

func (p *Parser) funcParse(input *nmcore.BlockStack) (nmcore.ExprNode, bool) {
	var object nmcore.ExprNode

	frame := *(input)

	//Check function have minimal items.
	if !(frame.Length() >= 4) {
		return object, false
	}

	//Check function is closed.
	if !(frame.Seek().Token == nmcore.TokenBlockEnd) {
		return object, false
	}

	//Check Open Token and Close Token is all pushed to stack.
	if !(p.bcounter.IsBlocked()) {
		return object, false
	}

	//Check function is starting with function symbol, Check function Open Token is existed in right place.
	//Clear buffer when Token block position is not correct.
	if !(frame.BlockList[0].Token == nmcore.TokenFunc && frame.BlockList[1].Token == nmcore.TokenBlockStart) {
		input.Clear()
		return object, false
	}

	//Is context string is Exist?
	if frame.BlockList[2].Token == nmcore.TokenString {
		funcData := frame.BlockList[2]               //Get context.
		splited := strings.Split(funcData.Item, " ") //Split string to parse function name & args

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

		object = nmcore.MakeFunctionNode(funcName, fnargs, []string{pFuncContext})

		//If inner-function is existed?
		if frame.BlockList[3].Token == nmcore.TokenFunc {
			var innerFuncStack nmcore.BlockStack
			innerFuncStack.AppendStack(frame, 3, frame.Length()-1)

			p.bcounter.DecCounter()
			obj, innerft := p.parseInnerBlock(innerFuncStack)
			if innerft {
				object.Insert(obj)
			}
		}

		input.Clear()

	}

	p.bcounter.DecCounter()

	return object, true
}

func (p *Parser) parseInnerBlock(input nmcore.BlockStack) ([]nmcore.ExprNode, bool) {
	var position = make([]int, 0)              //Position of deffunc
	var parseList []nmcore.BlockStack          //Parsed Innerfuncs
	var bkcounter = nmcore.MakeBraketCounter() //Counter for parse innerfunc braket
	var object []nmcore.ExprNode

	for idx, item := range input.BlockList {
		if item.Token == nmcore.TokenFunc {
			position = append(position, idx)
		}
	}

	//Return false when funcdef is not found.
	if len(position) == 0 || len(position) < 1 {
		return object, false
	}

	pointer := 0
	ipos := 0
	for pointer < len(input.BlockList) {
		switch input.BlockList[pointer].Token {
		case nmcore.TokenBlockStart:
			bkcounter.IncOpen()
		case nmcore.TokenBlockEnd:
			bkcounter.IncClose()
		}

		if bkcounter.IsBlocked() && (bkcounter.Open > 0 && bkcounter.Close > 0) {
			end := pointer + 1
			start := position[ipos]

			innerFuncStk := nmcore.NewBlockStack()
			innerFuncStk.AppendStack(input, start, end)

			parseList = append(parseList, innerFuncStk)
			bkcounter.Clear()

			for pi, t := range position {
				if t > pointer {
					ipos = pi
					pointer = position[pi]
					innerFuncStk.Clear()
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

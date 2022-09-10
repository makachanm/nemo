package nemomark

import (
	"strings"
)

type lexer struct {
}

type parser struct {
	func_buffer BlockStack
	bcounter    BraketCounter
}

func NewLexer() lexer {
	return lexer{}
}

func NewParser() parser {
	return parser{
		func_buffer: NewBlockStack(),
		bcounter:    MakeBraketCounter(),
	}
}

func return_string(s *string, idx int) string {
	return string((*s)[idx])
}

func add_block(b *[]Block, t Token, i string) []Block {
	return append((*b), Block{token: t, item: i})
}

func append_eol(b *[]Block) []Block {
	return append((*b), Block{token: TOKEN_EOL})
}

func (l *lexer) Tokenize(input string, tokenmap map[string]Token) []Block {
	var Blocks []Block
	var pointer int = 0
	var lengths int = len(input)

	for pointer < lengths {
		current_string := return_string(&input, pointer)
		token_val, key_exist := tokenmap[current_string]

		if key_exist {
			if token_val == TOKEN_IGNORE {
				iidx := pointer
				iis_break := (token_val == TOKEN_IGNORE)

				for iis_break {
					//Check idx is smaller than lengths. make break when string pointer is encountered eol.
					if (iidx + 1) == lengths {
						iidx++ //increase pointer
						break
					}
					t, r := tokenmap[return_string(&input, (iidx+1))]
					iis_break = !(r && t == TOKEN_IGNORE)
					iidx++
				}

				Blocks = add_block(&Blocks, TOKEN_STRING, input[(pointer+1):iidx])
				pointer = (iidx - pointer) + pointer + 1

			} else {
				//Tokenize current string.
				Blocks = add_block(&Blocks, token_val, "")
				pointer++
			}

		} else {
			//Make current text pointer to normal string.
			idx := pointer //String tokenizer pointer
			is_break := key_exist
			for !is_break { //loop until enconter another token.
				//Check idx is smaller than lengths. make break when string pointer is encountered eol.
				if (idx + 1) == lengths {
					idx++ //increase pointer
					break
				}
				//Check token.
				_, is_break = tokenmap[return_string(&input, (idx+1))]
				idx++ //increase pointer (not eol)
			}
			//Cut string and tokenize. [Strat of Stirng to End of String]
			Blocks = add_block(&Blocks, TOKEN_STRING, input[pointer:idx])
			pointer = (idx - pointer) + pointer
		}
	}
	//Append EOL token.
	Blocks = append_eol(&Blocks)

	return Blocks
}

func (p *parser) Parse(input *[]Block) ExprNode {
	var Blocks []Block = (*input)
	var Head ExprNode = MakeExprNode(TYPE_SECTION, nil) //Head of tree.
	var pointer int = 0
	var lengths int = len(Blocks)

	var Stack BlockStack = NewBlockStack() //Stack for StackParse.

	for pointer < lengths {
		Stack.block_push(Blocks[pointer])
		parsed := p.stack_parse(&Stack)
		if parsed.Context != nil || parsed.Func_context.Fucntion_name != "" {
			// Check Node's context is not a nil and Node is valid function node.
			Head.single_insert(parsed)
		}
		pointer++
	}

	return Head
}

func (p *parser) stack_parse(input *BlockStack) ExprNode {
	var object ExprNode

	switch input.seek().token {
	case TOKEN_STRING:
		if p.func_buffer.length() > 1 { //Check another block is exist in stack.
			p.func_buffer.block_push(input.block_pop())
			break
		}
		item := input.block_pop().item

		string_th := []string{item}
		object = MakeExprNode(TYPE_STRING, string_th)

	case TOKEN_BLOCK_START:
		open_token := input.block_pop()
		p.bcounter.inc_open()
		p.func_buffer.block_push(open_token)

	case TOKEN_BLOCK_END:
		close_token := input.block_pop()
		p.bcounter.inc_close()
		p.func_buffer.block_push(close_token)

		func_parsed, is_parsed := p.func_parse(&p.func_buffer)

		if is_parsed {
			object = func_parsed
		}

	default:
		p.func_buffer.block_push(input.block_pop())
		func_parsed, is_parsed := p.func_parse(&p.func_buffer)

		if is_parsed {
			object = func_parsed
		}
	}

	return object
}

func (p *parser) func_parse(input *BlockStack) (ExprNode, bool) {
	var object ExprNode

	frame := *(input)

	//Check function have minimal items.
	if !(frame.length() >= 4) {
		return object, false
	}

	//Check function is closed.
	if !(frame.seek().token == TOKEN_BLOCK_END) {
		return object, false
	}

	//Check open token and close token is all pushed to stack.
	if !(p.bcounter.is_blocked()) {
		return object, false
	}

	//Check function is starting with function symbol, Check function open token is exist in right place.
	//Clear buffer when token block position is not correct.
	if !(frame.block_list[0].token == TOKEN_FUNC && frame.block_list[1].token == TOKEN_BLOCK_START) {
		input.clear()
		return object, false
	}

	//Is context string is Exist?
	if frame.block_list[2].token == TOKEN_STRING {
		func_data := frame.block_list[2]              //Get context.
		splited := strings.Split(func_data.item, " ") //Split string to parse function name & args

		func_name := splited[0]                          //Get function name.
		p_func_context := strings.Join(splited[1:], " ") //Join remain strings to one context value.

		var fnargs map[string]string = make(map[string]string)

		if strings.ContainsAny(func_name, "(") && strings.ContainsAny(func_name, ")") {
			startpoint := strings.Index(func_name, "(")
			endpoint := strings.Index(func_name, ")")

			argstr := func_name[(startpoint + 1):endpoint]
			func_name = func_name[0:startpoint]

			args_splited := strings.Split(argstr, ",")
			for _, value := range args_splited {
				if strings.ContainsAny(value, "=") {
					args := strings.Split(value, "=")
					fnargs[args[0]] = strings.Join(args[1:], "")
				}
			}
		}

		object = MakeFunctionNode(func_name, fnargs, []string{p_func_context})

		//If inner-function is exist?
		if frame.block_list[3].token == TOKEN_FUNC {
			var inner_func_stack BlockStack
			inner_func_stack.append_stack(frame, 3, (frame.length() - 1))

			p.bcounter.dec_counter()
			obj, innerft := p.parse_inner_block(inner_func_stack)
			if innerft {
				object.insert(obj)
			}
		}

		input.clear()

	}

	p.bcounter.dec_counter()

	return object, true
}

func (p *parser) parse_inner_block(input BlockStack) ([]ExprNode, bool) {
	var position []int = make([]int, 0)               //Position of deffunc
	var parse_list []BlockStack                       //Parsed Innerfuncs
	var bkcounter BraketCounter = MakeBraketCounter() //Counter for parse innerfunc braket
	var object []ExprNode

	for idx, item := range input.block_list {
		if item.token == TOKEN_FUNC {
			position = append(position, idx)
		}
	}

	//Return false when funcdef is not found.
	if len(position) == 0 || len(position) < 1 {
		return object, false
	}

	pointer := 0
	ipos := 0
	for pointer < len(input.block_list) {
		switch input.block_list[pointer].token {
		case TOKEN_BLOCK_START:
			bkcounter.inc_open()
		case TOKEN_BLOCK_END:
			bkcounter.inc_close()
		}

		if bkcounter.is_blocked() && (bkcounter.open > 0 && bkcounter.close > 0) {
			end := pointer + 1
			start := position[ipos]

			inner_func_stk := NewBlockStack()
			inner_func_stk.append_stack(input, start, end)

			parse_list = append(parse_list, inner_func_stk)
			bkcounter.clear()

			for pi, t := range position {
				if t > pointer {
					ipos = pi
					pointer = position[pi]
					inner_func_stk.clear()
					break
				}
			}

		}

		pointer++
	}

	for _, context := range parse_list {
		pased, is_pased := p.func_parse(&context)
		if is_pased {
			object = append(object, pased)
		}
	}

	return object, true
}

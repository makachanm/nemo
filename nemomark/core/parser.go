package core

import (
	"strings"
)

type NLexer struct {
	internal_counter BraketCounter
}

type NParser struct {
	internal_counter BraketCounter
	func_stack       BlockStack
}

func NewNLexer() NLexer {
	return NLexer{
		internal_counter: MakeBraketCounter(),
	}
}

func NewNParser() NParser {
	return NParser{
		internal_counter: MakeBraketCounter(),
		func_stack:       NewBlockStack(),
	}
}

func (l *NLexer) Toknize(input string, tokenmap map[string]TokMapElement) []Block {
	var total_lengths uint64 = uint64(len(input))
	var pointer uint64 = 0

	var processed_blocks []Block

	for pointer < total_lengths {
		current_char := GetChar(&input, int(pointer))
		token_val, is_token := tokenmap[current_char]

		if is_token && token_val.Token != TokenIgnore {
			AppendSingleBlock(&processed_blocks, GenerateBlock(token_val.Token, token_val.MatchChar))
			pointer++
		} else if token_val.Token == TokenIgnore {
			var start_pos, end_pos uint64
			start_pos = pointer + 1
			current_pos := start_pos
			is_break := false

			for current_pos < total_lengths && !is_break {
				if tokenmap[GetChar(&input, int(current_pos))].Token == TokenIgnore {
					is_break = true
				}

				current_pos++
			}
			end_pos = current_pos

			AppendSingleBlock(&processed_blocks, GenerateBlock(TokenString, GetCharfromRange(&input, int(start_pos), int(end_pos-1))))
			pointer = uint64(end_pos)

		} else {
			var start_pos, end_pos uint64
			current_pos := pointer
			is_break := false

			start_pos = pointer
			for current_pos < total_lengths && !is_break {
				_, is_break = tokenmap[GetChar(&input, int(current_pos))]
				current_pos++
			}
			end_pos = current_pos - 1

			if total_lengths == (pointer + 1) { //if is last char of input?
				end_pos++
			}

			AppendSingleBlock(&processed_blocks, GenerateBlock(TokenString, GetCharfromRange(&input, int(start_pos), int(end_pos))))
			pointer = uint64(end_pos)

		}
	}

	AppendSingleBlock(&processed_blocks, GenerateBlock(TokenEol, ""))
	return processed_blocks
}

func (p *NParser) Parse(input []Block) ExprNode {
	var origin_node ExprNode = MakeExprNode(TypeSection, nil)
	var pointer uint64 = 0

	for input[pointer].Token != TokenEol {
		p.func_stack.BlockPush(input[pointer])

		switch p.func_stack.Seek().Token {
		case TokenString:
			if p.func_stack.Length() <= 1 {
				string_block := p.func_stack.BlockPop()
				origin_node.SingleInsert(MakeExprNode(TypeString, []string{string_block.Item}))
			}

		case TokenBlockStart:
			p.internal_counter.IncOpen()

		case TokenBlockEnd:
			p.internal_counter.IncClose()
			if p.internal_counter.IsBlocked() {
				rsnode := parseFuncStack(&p.func_stack)
				origin_node.SingleInsert(rsnode)
			}
		}

		pointer++

	}

	return origin_node
}

func parseFuncStack(input *BlockStack) ExprNode {
	var func_node ExprNode

	if !(input.Length() > 1) && input.Seek().Token == TokenString { //check it is singe string token
		func_node = MakeExprNode(TypeString, []string{input.BlockPop().Item})
		return func_node
	}

	stack_ctx := input.Clear()
	stack_ctx = stack_ctx[2:(len(stack_ctx) - 1)] //remove unused symbols

	func_node = parseFunc(stack_ctx[0].Item)
	RemoveElementBlockArray(&stack_ctx, 0)

	var inner_fn_stacks []BlockStack
	var sym_counter BraketCounter = MakeBraketCounter()
	var innerfns_pointer uint64 = 0
	var fstart, fend uint64 = 0, 0

	//parse inner blocks to pieces
	for len(stack_ctx) > 0 {
		switch stack_ctx[innerfns_pointer].Token {
		case TokenFunc:
			if !sym_counter.IsBlocked() {
				innerfns_pointer++
			} else {
				fstart = innerfns_pointer
				innerfns_pointer++
			}

		case TokenBlockStart:
			sym_counter.IncOpen()
			innerfns_pointer++

		case TokenBlockEnd:
			sym_counter.IncClose()
			if sym_counter.IsBlocked() {
				fend = innerfns_pointer

				nstack := NewBlockStack()
				nstack.AppendArray(&stack_ctx, int(fstart), int(fend))
				inner_fn_stacks = append(inner_fn_stacks, nstack)
				sym_counter.Clear()
				innerfns_pointer = 0
			} else if len(stack_ctx) <= 1 {
				RemoveElementBlockArray(&stack_ctx, 0)
			} else {
				innerfns_pointer++
			}

		default:
			if sym_counter.IsBlocked() || sym_counter.IsZero() {
				fend = innerfns_pointer

				nstack := NewBlockStack()
				nstack.AppendArray(&stack_ctx, int(fstart), int(fend))
				inner_fn_stacks = append(inner_fn_stacks, nstack)
				sym_counter.Clear()
				innerfns_pointer = 0

			} else {
				innerfns_pointer++
			}

		}
	}

	if len(inner_fn_stacks) > 0 {
		for _, ptaget_stack := range inner_fn_stacks {
			inp_node := parseFuncStack(&ptaget_stack)
			func_node.SingleInsert(inp_node)
		}
	}

	return func_node
}

func parseFunc(input string) ExprNode {
	splited_d := strings.Split(input, " ")
	fn_name := splited_d[0]
	fn_ctx := []string{strings.Join(splited_d[1:], " ")}
	fn_args := make(map[string]string)

	if strings.Contains(fn_name, "(") && strings.Contains(fn_name, ")") {
		split_point := (strings.Index(fn_name, "(") + 1)
		argsctx := fn_name[split_point:(len(fn_name) - 1)]
		qname := fn_name[:(split_point - 1)]

		ls := strings.Split(argsctx, ",")
		for _, ctxL := range ls {
			spt_arg := strings.Split(ctxL, "=")
			fn_args[spt_arg[0]] = strings.Join(spt_arg[1:], "=")
		}

		return MakeFunctionNode(qname, fn_args, fn_ctx)
	}

	return MakeFunctionNode(fn_name, fn_args, fn_ctx)
}

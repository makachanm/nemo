package nemomark

type Token int
type NodeType int

type Block struct {
	token Token
	item  string
}

type BlockStack struct {
	block_list []Block
}

type ExprNode struct {
	Node_type        NodeType
	Func_context     MarkdownFucntion
	Context          []string
	Has_child        bool
	Is_funciton_node bool
	Child            []ExprNode
}

type MarkdownFucntion struct {
	Fucntion_name string
	Args          map[string]string
	Context       []string
}

type BraketCounter struct {
	open  int
	close int
}

const (
	TOKEN_BLOCK_START Token = iota
	TOKEN_BLOCK_END
	TOKEN_EOL
	TOKEN_STRING
	TOKEN_FUNC
	TOKEN_IGNORE
)

const (
	TYPE_ORIGIN NodeType = iota
	TYPE_SECTION
	TYPE_FUNC
	TYPE_STRING
)

var TokenMap = map[string]Token{
	"$": TOKEN_FUNC,
	"[": TOKEN_BLOCK_START,
	"]": TOKEN_BLOCK_END,
	"`": TOKEN_IGNORE,
}

func NewBlockStack() BlockStack {
	return BlockStack{block_list: []Block{}}
}

func MakeExprNode(p NodeType, c []string) ExprNode {
	return ExprNode{
		Node_type: p,
		Context:   c,
	}
}

func MakeFunctionNode(fcname string, args map[string]string, context []string) ExprNode {
	mk_func := MarkdownFucntion{
		Fucntion_name: fcname,
		Args:          args,
		Context:       context,
	}
	exnode := MakeExprNode(TYPE_FUNC, nil)
	exnode.set_is_function(true)
	exnode.set_func_context(mk_func)

	return exnode
}

func MakeBraketCounter() BraketCounter {
	return BraketCounter{
		open:  0,
		close: 0,
	}
}

func (s *BlockStack) append_stack(si BlockStack, start int, end int) {
	s.block_list = append(s.block_list, si.block_list[start:end]...)
}

func (s *BlockStack) block_push(b Block) {
	s.block_list = append(s.block_list, b)
}

func (s *BlockStack) block_pop() Block {
	if len(s.block_list) == 0 {
		return Block{}
	}
	data := s.block_list[(len(s.block_list) - 1)]
	s.block_list = s.block_list[:(len(s.block_list) - 1)]
	return data
}

//BlockStack
func (s *BlockStack) clear() []Block {
	data := s.block_list
	s.block_list = nil
	return data
}

func (s *BlockStack) length() int {
	return len(s.block_list)
}

func (s *BlockStack) seek() Block {
	if len(s.block_list) == 0 {
		return Block{}
	}
	return s.block_list[(len(s.block_list) - 1)]
}

//ExprNode
func (e *ExprNode) single_insert(c ExprNode) {
	e.Child = append(e.Child, c)
	e.set_has_child(true)
}

func (e *ExprNode) insert(c []ExprNode) {
	e.Child = append(e.Child, c...)
	e.set_has_child(true)
}

func (e *ExprNode) set_has_child(b bool) {
	e.Has_child = b
}

func (e *ExprNode) set_is_function(b bool) {
	e.Is_funciton_node = b
}

func (e *ExprNode) set_func_context(c MarkdownFucntion) {
	e.Func_context = c
}

//BraketCounter

func (b *BraketCounter) inc_open() {
	b.open++
}

func (b *BraketCounter) inc_close() {
	b.close++
}

func (b *BraketCounter) dec_open() {
	b.open--
}

func (b *BraketCounter) dec_close() {
	b.close--
}

func (b *BraketCounter) dec_counter() {
	b.dec_open()
	b.dec_close()
}

func (b *BraketCounter) clear() {
	b.open = 0
	b.close = 0
}

func (b *BraketCounter) is_blocked() bool {
	return (b.open == b.close)
}

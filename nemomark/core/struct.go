package core

type Token int
type NodeType int

type TokMapElement struct {
	Token     Token
	MatchChar string
}

type Block struct {
	Token    Token
	Item     string
	StartPos int64
	EndPos   int64
}

type BlockStack struct {
	BlockList []Block
}

type ExprNode struct {
	NodeType       NodeType
	FuncContext    MarkdownFucntion
	Context        []string
	HasChild       bool
	IsFuncitonNode bool
	Child          []ExprNode
}

type MarkdownFucntion struct {
	FunctionName string
	Args         map[string]string
	Context      []string
}

type BraketCounter struct {
	Open  int
	Close int
}

const (
	TokenBlockStart Token = iota
	TokenBlockEnd
	TokenEol
	TokenString
	TokenFunc
	TokenIgnore
)

const (
	TypeSection NodeType = iota
	TypeFunc
	TypeString
)

var TokenMap = map[string]Token{
	"$": TokenFunc,
	"[": TokenBlockStart,
	"]": TokenBlockEnd,
	"`": TokenIgnore,
}

var NTokenMap = map[string]TokMapElement{
	"$": TokMapElement{Token: TokenFunc, MatchChar: "$"},
	"[": TokMapElement{Token: TokenBlockStart, MatchChar: "["},
	"]": TokMapElement{Token: TokenBlockEnd, MatchChar: "]"},
	"`": TokMapElement{Token: TokenIgnore, MatchChar: "`"},
}

func NewBlockStack() BlockStack {
	return BlockStack{BlockList: []Block{}}
}

func MakeExprNode(nodeType NodeType, context []string) ExprNode {
	return ExprNode{
		NodeType: nodeType,
		Context:  context,
	}
}

func MakeFunctionNode(functionName string, args map[string]string, context []string) ExprNode {
	mkFunc := MarkdownFucntion{
		FunctionName: functionName,
		Args:         args,
		Context:      context,
	}
	exnode := MakeExprNode(TypeFunc, nil)
	exnode.SetIsFunction(true)
	exnode.SetFuncContext(mkFunc)

	return exnode
}

func MakeBraketCounter() BraketCounter {
	return BraketCounter{
		Open:  0,
		Close: 0,
	}
}

func (s *BlockStack) AppendStack(si BlockStack, start int, end int) {
	s.BlockList = append(s.BlockList, si.BlockList[start:end]...)
}

func (s *BlockStack) AppendArray(bas *[]Block, start int, end int) {
	lsb := (*bas)[start:(end + 1)]
	s.BlockList = append(s.BlockList, lsb...)

	RemoveElementBlockArray(bas, end)
}

func (s *BlockStack) BlockPush(b Block) {
	s.BlockList = append(s.BlockList, b)
}

func (s *BlockStack) BlockPop() Block {
	if len(s.BlockList) == 0 {
		return Block{}
	}
	data := s.BlockList[(len(s.BlockList) - 1)]
	s.BlockList = s.BlockList[:(len(s.BlockList) - 1)]
	return data
}

// BlockStack

func (s *BlockStack) Clear() []Block {
	data := s.BlockList
	s.BlockList = nil
	return data
}

func (s *BlockStack) Length() int {
	return len(s.BlockList)
}

func (s *BlockStack) Seek() Block {
	if len(s.BlockList) == 0 {
		return Block{}
	}
	return s.BlockList[(len(s.BlockList) - 1)]
}

// ExprNode

func (e *ExprNode) SingleInsert(child ExprNode) {
	e.Child = append(e.Child, child)
	e.SetHasChild(true)
}

func (e *ExprNode) Insert(child []ExprNode) {
	e.Child = append(e.Child, child...)
	e.SetHasChild(true)
}

func (e *ExprNode) SetHasChild(hasChild bool) {
	e.HasChild = hasChild
}

func (e *ExprNode) SetIsFunction(isFunction bool) {
	e.IsFuncitonNode = isFunction
}

func (e *ExprNode) SetFuncContext(funcContext MarkdownFucntion) {
	e.FuncContext = funcContext
}

//BraketCounter

func (b *BraketCounter) IncOpen() {
	b.Open++
}

func (b *BraketCounter) IncClose() {
	b.Close++
}

func (b *BraketCounter) DecOpen() {
	b.Open--
}

func (b *BraketCounter) DecClose() {
	b.Close--
}

func (b *BraketCounter) DecCounter() {
	b.DecOpen()
	b.DecClose()
}

func (b *BraketCounter) Clear() {
	b.Open = 0
	b.Close = 0
}

func (b *BraketCounter) IsBlocked() bool {
	return (b.Open == b.Close) && ((b.Open != 0) && (b.Close != 0))
}

func (b *BraketCounter) IsZero() bool {
	return ((b.Open == 0) && (b.Close == 0))
}

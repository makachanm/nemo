package nemomark

type Token int
type NodeType int

type Block struct {
	token Token
	item  string
}

type BlockStack struct {
	blockList []Block
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
	FucntionName string
	Args         map[string]string
	Context      []string
}

type BraketCounter struct {
	open  int
	close int
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

func NewBlockStack() BlockStack {
	return BlockStack{blockList: []Block{}}
}

func MakeExprNode(p NodeType, c []string) ExprNode {
	return ExprNode{
		NodeType: p,
		Context:  c,
	}
}

func MakeFunctionNode(fcname string, args map[string]string, context []string) ExprNode {
	mkFunc := MarkdownFucntion{
		FucntionName: fcname,
		Args:         args,
		Context:      context,
	}
	exnode := MakeExprNode(TypeFunc, nil)
	exnode.setIsFunction(true)
	exnode.setFuncContext(mkFunc)

	return exnode
}

func MakeBraketCounter() BraketCounter {
	return BraketCounter{
		open:  0,
		close: 0,
	}
}

func (s *BlockStack) appendStack(si BlockStack, start int, end int) {
	s.blockList = append(s.blockList, si.blockList[start:end]...)
}

func (s *BlockStack) blockPush(b Block) {
	s.blockList = append(s.blockList, b)
}

func (s *BlockStack) blockPop() Block {
	if len(s.blockList) == 0 {
		return Block{}
	}
	data := s.blockList[(len(s.blockList) - 1)]
	s.blockList = s.blockList[:(len(s.blockList) - 1)]
	return data
}

// BlockStack
func (s *BlockStack) clear() []Block {
	data := s.blockList
	s.blockList = nil
	return data
}

func (s *BlockStack) length() int {
	return len(s.blockList)
}

func (s *BlockStack) seek() Block {
	if len(s.blockList) == 0 {
		return Block{}
	}
	return s.blockList[(len(s.blockList) - 1)]
}

// ExprNode
func (e *ExprNode) singleInsert(c ExprNode) {
	e.Child = append(e.Child, c)
	e.setHasChild(true)
}

func (e *ExprNode) insert(c []ExprNode) {
	e.Child = append(e.Child, c...)
	e.setHasChild(true)
}

func (e *ExprNode) setHasChild(b bool) {
	e.HasChild = b
}

func (e *ExprNode) setIsFunction(b bool) {
	e.IsFuncitonNode = b
}

func (e *ExprNode) setFuncContext(c MarkdownFucntion) {
	e.FuncContext = c
}

//BraketCounter

func (b *BraketCounter) incOpen() {
	b.open++
}

func (b *BraketCounter) incClose() {
	b.close++
}

func (b *BraketCounter) decOpen() {
	b.open--
}

func (b *BraketCounter) decClose() {
	b.close--
}

func (b *BraketCounter) decCounter() {
	b.decOpen()
	b.decClose()
}

func (b *BraketCounter) clear() {
	b.open = 0
	b.close = 0
}

func (b *BraketCounter) isBlocked() bool {
	return b.open == b.close
}

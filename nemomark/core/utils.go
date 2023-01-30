package core

func GetChar(s *string, idx int) string {
	return string((*s)[idx])
}

func GetCharfromRange(s *string, start int, end int) string {
	return string((*s)[start:end])
}

func GenerateBlock(tokentype Token, content string) Block {
	return Block{
		Token: tokentype,
		Item:  content,
	}
}

func AppendSingleBlock(bs *[]Block, input Block) {
	(*bs) = append((*bs), input)
}

func RemoveElementBlockArray(bs *[]Block, idx int) {
	(*bs) = (*bs)[(idx + 1):]
}

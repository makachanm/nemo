package nemomark

import (
	"strings"
)

type Renderer struct {
}

func NewRenderer() Renderer {
	return Renderer{}
}

func (rd *Renderer) itemRender(input ExprNode) string {
	var result strings.Builder

	if input.NodeType == TypeString {
		result.WriteString(RenderPlain(input.Context))
	} else {
		if handler, ok := MarkdownHandlers[input.FuncContext.FunctionName]; ok {
			renderTarget := input.FuncContext
			if input.HasChild {
				var childstrbuilder strings.Builder
				for _, childNode := range input.Child {
					childstrbuilder.WriteString(rd.itemRender(childNode))
				}
				renderTarget.Context = []string{childstrbuilder.String()}
			}
			result.WriteString(handler(renderTarget))
		}
	}

	return result.String()
}

func (rd *Renderer) Render(input ExprNode) string {
	if input.NodeType != TypeSection {
		return ""
	}

	var result string
	for _, node := range input.Child {
		rctx := rd.itemRender(node)
		result += rctx
	}

	return "<p>" + result + "</p>"
}

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
	var builded string
	var handler = MarkdownHandlers
	var stringbuilder strings.Builder

	if input.NodeType == TypeString {
		ctx := RenderPlain(input.Context)
		stringbuilder.WriteString(builded)
		stringbuilder.WriteString(ctx)

		builded = stringbuilder.String()
		stringbuilder.Reset()
	} else {
		handle, isFuncExist := handler[input.FuncContext.FucntionName]
		if isFuncExist {
			renderTarget := input.FuncContext
			childstr := ""
			if input.HasChild {
				var childstrbuilder strings.Builder
				for _, childNode := range input.Child {
					childctx := rd.itemRender(childNode)
					childstrbuilder.WriteString(childstr)
					childstrbuilder.WriteString(childctx)

					childstr = childstrbuilder.String()
					childstrbuilder.Reset()
				}
				renderTarget.Context = []string{childstr}
			}

			ctx := handle(renderTarget)
			stringbuilder.WriteString(builded)
			stringbuilder.WriteString(ctx)

			builded = stringbuilder.String()
			stringbuilder.Reset()
		}
	}

	return builded
}

func (rd *Renderer) render(input ExprNode) string {
	var builded string
	//var handler = Markdown_Handlers
	var stringbuilder strings.Builder

	originNode := input
	if !(originNode.NodeType == TypeSection) {
		return ""
	}

	originChildNode := originNode.Child

	for _, node := range originChildNode {
		ctx := rd.itemRender(node)
		stringbuilder.WriteString(builded)
		stringbuilder.WriteString(ctx)

		builded = stringbuilder.String()
		stringbuilder.Reset()
	}

	taggedResult := "<p>" + builded + "</p>"
	return taggedResult
}

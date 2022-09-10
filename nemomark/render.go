package nemomark

import (
	"strings"
)

type renderer struct {
}

func NewRenderer() renderer {
	return renderer{}
}

func (rd *renderer) item_render(input ExprNode) string {
	var builded string
	var handler = Markdown_Handlers
	var stringbuilder strings.Builder

	if input.Node_type == TYPE_STRING {
		ctx := RenderPlain(input.Context)
		stringbuilder.WriteString(builded)
		stringbuilder.WriteString(ctx)

		builded = stringbuilder.String()
		stringbuilder.Reset()
	} else {
		handle, is_func_exist := handler[input.Func_context.Fucntion_name]
		if is_func_exist {
			render_target := input.Func_context
			childstr := ""
			if input.Has_child {
				var childstrbuilder strings.Builder
				for _, child_node := range input.Child {
					childctx := rd.item_render(child_node)
					childstrbuilder.WriteString(childstr)
					childstrbuilder.WriteString(childctx)

					childstr = childstrbuilder.String()
					childstrbuilder.Reset()
				}
				render_target.Context = []string{childstr}
			}

			ctx := handle(render_target)
			stringbuilder.WriteString(builded)
			stringbuilder.WriteString(ctx)

			builded = stringbuilder.String()
			stringbuilder.Reset()
		}
	}

	return builded
}

func (rd *renderer) render(input ExprNode) string {
	var builded string
	//var handler = Markdown_Handlers
	var stringbuilder strings.Builder

	origin_node := input
	if !(origin_node.Node_type == TYPE_SECTION) {
		return ""
	}

	origin_child_node := origin_node.Child

	for _, node := range origin_child_node {
		ctx := rd.item_render(node)
		stringbuilder.WriteString(builded)
		stringbuilder.WriteString(ctx)

		builded = stringbuilder.String()
		stringbuilder.Reset()
	}

	tagged_result := "<p>" + builded + "</p>"
	return tagged_result
}

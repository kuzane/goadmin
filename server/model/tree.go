package model

type TreeNode struct {
    Label    string     `json:"label"`
    Value    string     `json:"value"`
    Children []TreeNode `json:"children,omitempty"`
}

// 菜单树
func GenerateTree(endpoints []*Endpoint) []TreeNode {
    tree := make(map[string]map[string][]TreeNode)

    // 构建树状结构
    for _, ep := range endpoints {
        if tree[ep.Module] == nil {
            tree[ep.Module] = make(map[string][]TreeNode)
        }

        node := TreeNode{
            Value: ep.Identity, // 可以根据需要设置 ID
            Label: ep.Remark,
        }

        tree[ep.Module][ep.Kind] = append(tree[ep.Module][ep.Kind], node)
    }

    // 转换为最终结果格式
    var result []TreeNode
    for module, kinds := range tree {
        moduleNode := TreeNode{
            Value: module,
            Label: module,
        }

        for kind, nodes := range kinds {
            kindNode := TreeNode{
                Value:    kind,
                Label:    kind,
                Children: nodes,
            }
            moduleNode.Children = append(moduleNode.Children, kindNode)
        }

        result = append(result, moduleNode)
    }

    return result
}

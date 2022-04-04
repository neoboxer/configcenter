package core

import "strings"

type trie struct {
	path   string           // full URL path
	part   string           // part of URL path
	son    map[string]*trie // trie's son nodes
	isWild bool             // mark wild card
}

func (t *trie) insert(path string) {
	parts := parsePath(path)
	root := t
	for _, part := range parts {
		if root.son[part] == nil {
			root.son[part] = &trie{
				part:   part,
				son:    make(map[string]*trie),
				isWild: part[0] == '*',
			}
		}
		root = root.son[part]
	}
	root.path = path
}

func (t *trie) search(path string) (*trie, map[string]string) {
	parts := parsePath(path)
	params := map[string]string{}
	root := t
	for i, part := range parts {
		var temp string
		for _, node := range root.son {
			if node.isWild {
				params[node.part[1:]] = parts[i]
			}
			if node.part == part || node.isWild {
				temp = node.part
			}
		}
		// path not match in trie
		if temp == "" {
			return nil, nil
		}
		if temp[0] == '*' {
			return root.son[temp], params
		}
		root = root.son[temp]
	}
	return root, params
}

func parsePath(path string) []string {
	res := make([]string, 0)
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part != "" {
			res = append(res, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return res
}

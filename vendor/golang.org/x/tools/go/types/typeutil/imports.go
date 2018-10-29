



package typeutil

import "go/types"








func Dependencies(pkgs ...*types.Package) []*types.Package {
	var result []*types.Package
	seen := make(map[*types.Package]bool)
	var visit func(pkgs []*types.Package)
	visit = func(pkgs []*types.Package) {
		for _, p := range pkgs {
			if !seen[p] {
				seen[p] = true
				visit(p.Imports())
				result = append(result, p)
			}
		}
	}
	visit(pkgs)
	return result
}

// Code generated by mknode.go. DO NOT EDIT.

package ir

import "fmt"

func (n *AddStringExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *AddStringExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.List = c.List.Copy()
	return &c
}
func (n *AddStringExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDoList(n.List, err, do)
	return err
}
func (n *AddStringExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	editList(n.List, edit)
}

func (n *AddrExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *AddrExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *AddrExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Alloc, err, do)
	return err
}
func (n *AddrExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Alloc = maybeEdit(n.Alloc, edit)
}

func (n *ArrayType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ArrayType) copy() Node {
	c := *n
	return &c
}
func (n *ArrayType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDo(n.Len, err, do)
	err = maybeDo(n.Elem, err, do)
	return err
}
func (n *ArrayType) editChildren(edit func(Node) Node) {
	n.Len = maybeEdit(n.Len, edit)
	n.Elem = maybeEdit(n.Elem, edit)
}

func (n *AssignListStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *AssignListStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Lhs = c.Lhs.Copy()
	c.Rhs = c.Rhs.Copy()
	return &c
}
func (n *AssignListStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDoList(n.Lhs, err, do)
	err = maybeDoList(n.Rhs, err, do)
	return err
}
func (n *AssignListStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	editList(n.Lhs, edit)
	editList(n.Rhs, edit)
}

func (n *AssignOpStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *AssignOpStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *AssignOpStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Y, err, do)
	return err
}
func (n *AssignOpStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Y = maybeEdit(n.Y, edit)
}

func (n *AssignStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *AssignStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *AssignStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Y, err, do)
	return err
}
func (n *AssignStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Y = maybeEdit(n.Y, edit)
}

func (n *BasicLit) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *BasicLit) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *BasicLit) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *BasicLit) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *BinaryExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *BinaryExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *BinaryExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Y, err, do)
	return err
}
func (n *BinaryExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Y = maybeEdit(n.Y, edit)
}

func (n *BlockStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *BlockStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.List = c.List.Copy()
	return &c
}
func (n *BlockStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDoList(n.List, err, do)
	return err
}
func (n *BlockStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	editList(n.List, edit)
}

func (n *BranchStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *BranchStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *BranchStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *BranchStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *CallExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *CallExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Args = c.Args.Copy()
	c.Rargs = c.Rargs.Copy()
	c.Body = c.Body.Copy()
	return &c
}
func (n *CallExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDoList(n.Args, err, do)
	err = maybeDoList(n.Rargs, err, do)
	err = maybeDoList(n.Body, err, do)
	return err
}
func (n *CallExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	editList(n.Args, edit)
	editList(n.Rargs, edit)
	editList(n.Body, edit)
}

func (n *CallPartExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *CallPartExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *CallPartExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	return err
}
func (n *CallPartExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
}

func (n *CaseStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *CaseStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.List = c.List.Copy()
	c.Body = c.Body.Copy()
	return &c
}
func (n *CaseStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Var, err, do)
	err = maybeDoList(n.List, err, do)
	err = maybeDoList(n.Body, err, do)
	return err
}
func (n *CaseStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Var = maybeEdit(n.Var, edit)
	editList(n.List, edit)
	editList(n.Body, edit)
}

func (n *ChanType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ChanType) copy() Node {
	c := *n
	return &c
}
func (n *ChanType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDo(n.Elem, err, do)
	return err
}
func (n *ChanType) editChildren(edit func(Node) Node) {
	n.Elem = maybeEdit(n.Elem, edit)
}

func (n *ClosureExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ClosureExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *ClosureExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *ClosureExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *ClosureReadExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ClosureReadExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *ClosureReadExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *ClosureReadExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *CommStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *CommStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Body = c.Body.Copy()
	return &c
}
func (n *CommStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Comm, err, do)
	err = maybeDoList(n.Body, err, do)
	return err
}
func (n *CommStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Comm = maybeEdit(n.Comm, edit)
	editList(n.Body, edit)
}

func (n *CompLitExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *CompLitExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.List = c.List.Copy()
	return &c
}
func (n *CompLitExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Ntype, err, do)
	err = maybeDoList(n.List, err, do)
	return err
}
func (n *CompLitExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Ntype = toNtype(maybeEdit(n.Ntype, edit))
	editList(n.List, edit)
}

func (n *ConstExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ConstExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *ConstExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *ConstExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *ConvExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ConvExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *ConvExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	return err
}
func (n *ConvExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
}

func (n *Decl) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *Decl) copy() Node {
	c := *n
	return &c
}
func (n *Decl) doChildren(do func(Node) error) error {
	var err error
	err = maybeDo(n.X, err, do)
	return err
}
func (n *Decl) editChildren(edit func(Node) Node) {
	n.X = maybeEdit(n.X, edit)
}

func (n *ForStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ForStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Late = c.Late.Copy()
	c.Body = c.Body.Copy()
	return &c
}
func (n *ForStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Cond, err, do)
	err = maybeDoList(n.Late, err, do)
	err = maybeDo(n.Post, err, do)
	err = maybeDoList(n.Body, err, do)
	return err
}
func (n *ForStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Cond = maybeEdit(n.Cond, edit)
	editList(n.Late, edit)
	n.Post = maybeEdit(n.Post, edit)
	editList(n.Body, edit)
}

func (n *Func) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *Func) copy() Node {
	c := *n
	c.Body = c.Body.Copy()
	return &c
}
func (n *Func) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.Body, err, do)
	return err
}
func (n *Func) editChildren(edit func(Node) Node) {
	editList(n.Body, edit)
}

func (n *FuncType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *FuncType) copy() Node {
	c := *n
	if c.Recv != nil {
		c.Recv = c.Recv.copy()
	}
	c.Params = copyFields(c.Params)
	c.Results = copyFields(c.Results)
	return &c
}
func (n *FuncType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoField(n.Recv, err, do)
	err = maybeDoFields(n.Params, err, do)
	err = maybeDoFields(n.Results, err, do)
	return err
}
func (n *FuncType) editChildren(edit func(Node) Node) {
	editField(n.Recv, edit)
	editFields(n.Params, edit)
	editFields(n.Results, edit)
}

func (n *GoDeferStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *GoDeferStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *GoDeferStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Call, err, do)
	return err
}
func (n *GoDeferStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Call = maybeEdit(n.Call, edit)
}

func (n *Ident) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *Ident) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *Ident) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *Ident) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *IfStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *IfStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Body = c.Body.Copy()
	c.Else = c.Else.Copy()
	return &c
}
func (n *IfStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Cond, err, do)
	err = maybeDoList(n.Body, err, do)
	err = maybeDoList(n.Else, err, do)
	return err
}
func (n *IfStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Cond = maybeEdit(n.Cond, edit)
	editList(n.Body, edit)
	editList(n.Else, edit)
}

func (n *IndexExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *IndexExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *IndexExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Index, err, do)
	return err
}
func (n *IndexExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Index = maybeEdit(n.Index, edit)
}

func (n *InlineMarkStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *InlineMarkStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *InlineMarkStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *InlineMarkStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *InlinedCallExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *InlinedCallExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Body = c.Body.Copy()
	c.ReturnVars = c.ReturnVars.Copy()
	return &c
}
func (n *InlinedCallExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDoList(n.Body, err, do)
	err = maybeDoList(n.ReturnVars, err, do)
	return err
}
func (n *InlinedCallExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	editList(n.Body, edit)
	editList(n.ReturnVars, edit)
}

func (n *InterfaceType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *InterfaceType) copy() Node {
	c := *n
	c.Methods = copyFields(c.Methods)
	return &c
}
func (n *InterfaceType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoFields(n.Methods, err, do)
	return err
}
func (n *InterfaceType) editChildren(edit func(Node) Node) {
	editFields(n.Methods, edit)
}

func (n *KeyExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *KeyExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *KeyExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Key, err, do)
	err = maybeDo(n.Value, err, do)
	return err
}
func (n *KeyExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Key = maybeEdit(n.Key, edit)
	n.Value = maybeEdit(n.Value, edit)
}

func (n *LabelStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *LabelStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *LabelStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *LabelStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *LogicalExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *LogicalExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *LogicalExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Y, err, do)
	return err
}
func (n *LogicalExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Y = maybeEdit(n.Y, edit)
}

func (n *MakeExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *MakeExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *MakeExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Len, err, do)
	err = maybeDo(n.Cap, err, do)
	return err
}
func (n *MakeExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Len = maybeEdit(n.Len, edit)
	n.Cap = maybeEdit(n.Cap, edit)
}

func (n *MapType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *MapType) copy() Node {
	c := *n
	return &c
}
func (n *MapType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDo(n.Key, err, do)
	err = maybeDo(n.Elem, err, do)
	return err
}
func (n *MapType) editChildren(edit func(Node) Node) {
	n.Key = maybeEdit(n.Key, edit)
	n.Elem = maybeEdit(n.Elem, edit)
}

func (n *MethodExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *MethodExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *MethodExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *MethodExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *Name) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *Name) copy() Node                    { panic("Name.copy") }
func (n *Name) doChildren(do func(Node) error) error {
	var err error
	return err
}
func (n *Name) editChildren(edit func(Node) Node) {
}

func (n *NameOffsetExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *NameOffsetExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *NameOffsetExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *NameOffsetExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *NilExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *NilExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *NilExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *NilExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *ParenExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ParenExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *ParenExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	return err
}
func (n *ParenExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
}

func (n *PkgName) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *PkgName) copy() Node {
	c := *n
	return &c
}
func (n *PkgName) doChildren(do func(Node) error) error {
	var err error
	return err
}
func (n *PkgName) editChildren(edit func(Node) Node) {
}

func (n *RangeStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *RangeStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Body = c.Body.Copy()
	return &c
}
func (n *RangeStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Key, err, do)
	err = maybeDo(n.Value, err, do)
	err = maybeDoList(n.Body, err, do)
	return err
}
func (n *RangeStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Key = maybeEdit(n.Key, edit)
	n.Value = maybeEdit(n.Value, edit)
	editList(n.Body, edit)
}

func (n *ResultExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ResultExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *ResultExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	return err
}
func (n *ResultExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
}

func (n *ReturnStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *ReturnStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Results = c.Results.Copy()
	return &c
}
func (n *ReturnStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDoList(n.Results, err, do)
	return err
}
func (n *ReturnStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	editList(n.Results, edit)
}

func (n *SelectStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SelectStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Cases = copyComms(c.Cases)
	c.Compiled = c.Compiled.Copy()
	return &c
}
func (n *SelectStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDoComms(n.Cases, err, do)
	err = maybeDoList(n.Compiled, err, do)
	return err
}
func (n *SelectStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	editComms(n.Cases, edit)
	editList(n.Compiled, edit)
}

func (n *SelectorExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SelectorExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *SelectorExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	return err
}
func (n *SelectorExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
}

func (n *SendStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SendStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *SendStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Chan, err, do)
	err = maybeDo(n.Value, err, do)
	return err
}
func (n *SendStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Chan = maybeEdit(n.Chan, edit)
	n.Value = maybeEdit(n.Value, edit)
}

func (n *SliceExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SliceExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *SliceExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Low, err, do)
	err = maybeDo(n.High, err, do)
	err = maybeDo(n.Max, err, do)
	return err
}
func (n *SliceExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Low = maybeEdit(n.Low, edit)
	n.High = maybeEdit(n.High, edit)
	n.Max = maybeEdit(n.Max, edit)
}

func (n *SliceHeaderExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SliceHeaderExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *SliceHeaderExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Ptr, err, do)
	err = maybeDo(n.Len, err, do)
	err = maybeDo(n.Cap, err, do)
	return err
}
func (n *SliceHeaderExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Ptr = maybeEdit(n.Ptr, edit)
	n.Len = maybeEdit(n.Len, edit)
	n.Cap = maybeEdit(n.Cap, edit)
}

func (n *SliceType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SliceType) copy() Node {
	c := *n
	return &c
}
func (n *SliceType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDo(n.Elem, err, do)
	return err
}
func (n *SliceType) editChildren(edit func(Node) Node) {
	n.Elem = maybeEdit(n.Elem, edit)
}

func (n *StarExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *StarExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *StarExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	return err
}
func (n *StarExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
}

func (n *StructKeyExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *StructKeyExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *StructKeyExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Value, err, do)
	return err
}
func (n *StructKeyExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Value = maybeEdit(n.Value, edit)
}

func (n *StructType) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *StructType) copy() Node {
	c := *n
	c.Fields = copyFields(c.Fields)
	return &c
}
func (n *StructType) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoFields(n.Fields, err, do)
	return err
}
func (n *StructType) editChildren(edit func(Node) Node) {
	editFields(n.Fields, edit)
}

func (n *SwitchStmt) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *SwitchStmt) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Cases = copyCases(c.Cases)
	c.Compiled = c.Compiled.Copy()
	return &c
}
func (n *SwitchStmt) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.Tag, err, do)
	err = maybeDoCases(n.Cases, err, do)
	err = maybeDoList(n.Compiled, err, do)
	return err
}
func (n *SwitchStmt) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.Tag = maybeEdit(n.Tag, edit)
	editCases(n.Cases, edit)
	editList(n.Compiled, edit)
}

func (n *TypeAssertExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *TypeAssertExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	c.Itab = c.Itab.Copy()
	return &c
}
func (n *TypeAssertExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	err = maybeDo(n.Ntype, err, do)
	err = maybeDoList(n.Itab, err, do)
	return err
}
func (n *TypeAssertExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
	n.Ntype = maybeEdit(n.Ntype, edit)
	editList(n.Itab, edit)
}

func (n *TypeSwitchGuard) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *TypeSwitchGuard) copy() Node {
	c := *n
	return &c
}
func (n *TypeSwitchGuard) doChildren(do func(Node) error) error {
	var err error
	if n.Tag != nil {
		err = maybeDo(n.Tag, err, do)
	}
	err = maybeDo(n.X, err, do)
	return err
}
func (n *TypeSwitchGuard) editChildren(edit func(Node) Node) {
	if n.Tag != nil {
		n.Tag = edit(n.Tag).(*Ident)
	}
	n.X = maybeEdit(n.X, edit)
}

func (n *UnaryExpr) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *UnaryExpr) copy() Node {
	c := *n
	c.init = c.init.Copy()
	return &c
}
func (n *UnaryExpr) doChildren(do func(Node) error) error {
	var err error
	err = maybeDoList(n.init, err, do)
	err = maybeDo(n.X, err, do)
	return err
}
func (n *UnaryExpr) editChildren(edit func(Node) Node) {
	editList(n.init, edit)
	n.X = maybeEdit(n.X, edit)
}

func (n *typeNode) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }
func (n *typeNode) copy() Node {
	c := *n
	return &c
}
func (n *typeNode) doChildren(do func(Node) error) error {
	var err error
	return err
}
func (n *typeNode) editChildren(edit func(Node) Node) {
}

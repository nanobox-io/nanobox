package mist

import "sort"

// interfaces
type (

	//
	subscriptions interface {
		Add([]string)
		Remove([]string)
		Match([]string) bool
		ToSlice() [][]string
	}
)

type (

	// Node ...
	Node struct {
		branches map[string]*Node
		leaves   map[string]struct{}
	}
)

//
func newNode() (node *Node) {

	node = &Node{
		branches: map[string]*Node{},
		leaves:   map[string]struct{}{},
	}

	//
	return
}

// Add sorts the keys and then attempts to add them
func (node *Node) Add(keys []string) {

	//
	if len(keys) == 0 {
		return
	}

	sort.Strings(keys)
	node.add(keys)
}

// add ...
func (node *Node) add(keys []string) {

	// if there is only one key remaining we are at the end of the chain, so a leaf
	// is created
	if len(keys) == 1 {
		node.leaves[keys[0]] = struct{}{}
		return
	}

	// see if there is already a branch for the first key. if not create a new one
	// and add it; if a branch already exists we simply use it and continue
	branch, ok := node.branches[keys[0]]
	if !ok {
		branch = newNode()
		node.branches[keys[0]] = branch
	}

	// for the current branch (new or existing) continue to the next set of keys
	// adding them as branches until there is only one left, at which point a leaf
	// is created (above)
	branch.add(keys[1:])
}

// Remove sorts the keys and then attempts to remove them
func (node *Node) Remove(keys []string) {

	//
	if len(keys) == 0 {
		return
	}

	sort.Strings(keys)
	node.remove(keys)
}

// remove ...
func (node *Node) remove(keys []string) {

	// if there is only one key remaining we are at the end of the chain and need
	// to remove just the leaf
	if len(keys) == 1 {
		delete(node.leaves, keys[0])
		return
	}

	// see if a branch for the first key exists; if a branch exists we need to
	// recurse down the branch until we reach the end...
	branch, ok := node.branches[keys[0]]
	if ok {

		// continue key by key until we reach a leaf at which point its removed (above)
		branch.remove(keys[1:])

		// NOTE: this cleanup is a little inefficient and probably not needed unless
		// mist starts using a ton of memory
		//
		// once we reach the end of the line, if there are no more leaves or branch
		// on this branch, we can remove the branch
		// if len(branch.leaves) == 0 && len(branch.branches) == 0 {
		// 	delete(node.branches, keys[0])
		// }
	}
}

// Match sorts the keys and then attempts to find a match
func (node *Node) Match(keys []string) bool {
	sort.Strings(keys)
	return node.match(keys)
}

// ​match ...
func (node *Node) match(keys []string) bool {

	//
	if len(keys) == 0 {
		return false
	}

	// iterate through each key looking for a leaf, if found it's a match
	for _, key := range keys {
		if _, ok := node.leaves[key]; ok {
			return true
		}
	}

	// see if a branch for the first key exists; if no branch exists we need to
	// try and find a match for the next key, continuing down the chain until we
	// find a leaf (above)
	branch, ok := node.branches[keys[0]]
	if !ok {
		return node.match(keys[1:])
	}

	// if a branch does exist we down the branch until we find a leaf (above)
	return branch.match(keys[1:])
}

// ToSlice recurses down an entire node returning a list of all branches and leaves
// as a slice of slices
func (node *Node) ToSlice() (list [][]string) {

	// iterate through each leaf appending it as a slice to the list of keys
	for leaf := range node.leaves {
		list = append(list, []string{leaf})
	}

	// iterate through each branch getting its list of branches and appending those
	// to the list
	for branch, node := range node.branches {

		// get the current nodes slice of branches and leaves
		nodeSlice := node.ToSlice()

		// for each branch in the nodes list apppend the key to that key
		for _, nodeKey := range nodeSlice {
			list = append(list, append(nodeKey, branch))
		}
	}

	// sort each list
	for _, l := range list {
		sort.Strings(l)
	}

	return
}

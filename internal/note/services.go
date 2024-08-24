package note

import "slices"

// FindUids returns all Uids found in given text.
func FindUids(text string) []string {
	uidRe := getUidRegexp()
	return uidRe.FindAllString(text, -1)
}

// FindReferences returns map of references between Notes. Key is Uid of a Note
// in which other Uids will be looked for. If any are found, they populate
// slice returned in map's value.
func FindReferences(repository INoteRepository) ReferenceMap {
	uids, _ := repository.List()
	var refersTo = make(ReferenceMap)

	for _, uid := range uids {
		nt, _ := repository.Get(uid)
		currentNoteRefersTo := FindUids(nt.Body)
		slices.Sort(currentNoteRefersTo)
		currentNoteRefersTo = slices.Compact(currentNoteRefersTo)

		refersTo[uid] = currentNoteRefersTo
	}

	return refersTo
}

// ReverseReferences inverts given references map by swapping keys with values.
// If values length >1, then many keys are created.
func ReverseReferences(refersTo ReferenceMap) ReferenceMap {
	var referredBy = make(ReferenceMap)

	for root, refs := range refersTo {
		for _, ref := range refs {
			referredByRef := append(referredBy[ref], root)
			referredBy[ref] = referredByRef
		}
	}

	for root, refs := range referredBy {
		slices.Sort(refs)
		refs = slices.Compact(refs)
		referredBy[root] = refs
	}

	return referredBy
}

package explorer

import "github.com/EnotInc/Bard/internal/enums/keys"

func (ex *Explorer) beginDeletion() {
	ex.action = deleting
}

func (ex *Explorer) handleDeletion(key rune) {
	switch key {
	case keys.Enter, 'y', 'Y':
		ex.delFileWithCallback()
		ex.update = true
	case keys.Esc, 'n', 'N':
	default:
	}
	ex.action = none
}

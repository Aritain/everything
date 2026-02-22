package entry

import (
	"everything/common"
	"everything/models"
	e "everything/models/entry"
)

func EntryCreationStart(userID int64, ec *[]e.Entry) (mr models.ModuleResponse) {
	*ec = append(*ec, e.Entry{UserID: userID})
	mr.Text = "Reminder name?"
	mr.Keyboard = common.CompileCancelKeyboard()
	return mr
}

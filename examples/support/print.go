package support

import (
	"fmt"

	"github.com/iaintshine/basecrm-go/basecrm"

	tm "github.com/buger/goterm"
)

func accountView(a *basecrm.Account) string {
	return fmt.Sprintf("%d\t%s\t%s\t%s\n", a.Id, a.Name, a.Phone, a.Plan)
}

func userView(u *basecrm.User) string {
	return fmt.Sprintf("%d\t%s\t%s\t%s\n", u.Id, u.Name, u.Email, u.Role)
}

func PrintHeader(header string) {
	tm.Println(tm.Color(header, tm.GREEN))
}

func PrintWhoami(me *basecrm.User) {
	PrintHeader("[WHOAMI]")
	table := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(table, "Id\tName\tEmail\tRole\n")
	fmt.Fprintf(table, userView(me))
	table.String()
	tm.Println(table)
}

func PrintTeam(team []*basecrm.User) {
	PrintHeader("[TEAM]")
	table := tm.NewTable(0, 10, 5, ' ', 0)
	table.String()
	fmt.Fprintf(table, "Id\tName\tEmail\tRole\n")
	for _, member := range team {
		fmt.Fprintf(table, userView(member))
	}
	table.String()
	tm.Println(table)
}

func PrintCompany(company *basecrm.Account) {
	PrintHeader("[COMPANY]")
	table := tm.NewTable(0, 10, 5, ' ', 0)
	table.String()
	fmt.Fprintf(table, "Id\tName\tPhone\tPlan\n")
	fmt.Fprintf(table, accountView(company))
	table.String()
	tm.Println(table)
}

func Flush() {
	tm.Flush()
}

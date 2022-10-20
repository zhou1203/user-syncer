package types

type User struct {
	ID      string `json:"id" db:"USER_ID"`
	LoginNo string `db:"LOGIN_NO"`
	Status  int    `json:"status" db:"APP_ACCT_STATUS"`
	Name    string `json:"username" db:"USER_NAME"`
	Email   string `json:"email" db:"EMAIL"`
	OrgID   string `json:"orgId" db:"ORG_ID"`
	Mobile  string `json:"mobile" db:"MOBILE"`
	Source  string `json:"-"`
}

package types

type User struct {
	ID      int64  `json:"id" db:"USER_ID"`
	LoginNo string `db:"LOGIN_NO"`
	Status  int    `json:"status" db:"APP_ACCT_STATUS"` // http response: 0 valid 1 invalid, db: 1 valid 0 invalid
	Name    string `json:"username" db:"USER_NAME"`
	Email   string `json:"email" db:"EMAIL"`
	OrgID   string `json:"orgCode" db:"ORG_ID"`
	Mobile  string `json:"mobile" db:"MOBILE"`
	Source  string `json:"-"`
}

type Org struct {
	ID          string `json:"orgCode" db:"ORG_ID"`
	OrgName     string `json:"fullName" db:"ORG_NAME"`
	ParentOrgID string `json:"parentCode" db:"PARENT_ORG_ID"`
	Status      int    `json:"status"` // 0 valid, 1 invalid
}

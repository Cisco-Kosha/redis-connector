package models

type Account struct {
	OrganisationID   int64       `json:"organisation_id,omitempty"`
	OrganisationName string      `json:"organisation_name,omitempty"`
	AccountID        int         `json:"account_id,omitempty"`
	AccountName      string      `json:"account_name,omitempty"`
	AccountDomain    string      `json:"account_domain,omitempty"`
	BundleID         interface{} `json:"bundle_id,omitempty"`
	HipaaCompliant   bool        `json:"hipaa_compliant,omitempty"`
	TotalAgents      struct {
		FullTime      int `json:"full_time,omitempty"`
		Occasional    int `json:"occasional,omitempty"`
		FieldService  int `json:"field_service,omitempty"`
		Collaborators int `json:"collaborators,omitempty"`
	} `json:"total_agents,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
	DataCenter string `json:"data_center,omitempty"`
	TierType   string `json:"tier_type,omitempty"`
	Address    struct {
		Country    string `json:"country,omitempty"`
		State      string `json:"state,omitempty"`
		City       string `json:"city,omitempty"`
		Street     string `json:"street,omitempty"`
		Postalcode string `json:"postalcode,omitempty"`
	} `json:"address,omitempty"`
	ContactPerson struct {
		Firstname string `json:"firstname,omitempty"`
		Lastname  string `json:"lastname,omitempty"`
		Email     string `json:"email,omitempty"`
	} `json:"contact_person,omitempty"`
}

type Specification struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Db        string `json:"database,omitempty"`
	RedisHost string `json:"redis_host,omitempty"`
	RedisPort string `json:"redis_port,omitempty"`
}

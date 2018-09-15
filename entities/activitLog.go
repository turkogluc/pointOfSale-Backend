package entities

type ActivityLogItem struct {
	User	string 	`json:"user"`
	Date	int		`json:"date"`
	ActivityType	string	`json:"activityType"`
	Title		string	`json:"title"`
	Description	string	`json:"description"`
	Detail	string `json:"detail"`
}


type ActivityLogs struct {
	Count	int	`json:"count"`
	Items	[]*ActivityLogItem	`json:"items"`
}
package form

type UpdateUserBalance struct {
	Username   string `json:"username"`
	AmountUsdt string `json:"amountUsdt"`
}

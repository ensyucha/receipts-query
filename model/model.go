package model

type User struct {
	UserId int `json:"userid"`
	Username string `json:"username"`
	NickName string `json:"nickname"`
	Password string `json:"password"`
	Usage int `json:"usage"`
	Total int `json:"total"`
}

type System struct {
	Password string `json:"password"`
	UnusedUsage int `json:"unusedusage"`
	ApiCode string `json:"apicode"`
}

type Query struct {
	Fpdm string `json:"fpdm"`
	Fphm string `json:"fphm"`
	Kprq string `json:"kprq"`
	Je   string `json:"je"`
	Jym  string `json:"jym"`
}

type QueryArray struct {
	QueryArray []*Query `json:"queryArray"`
}

type QueryResult struct {
	RespCode string `json:"respCode"`
	RespMsg string `json:"respMsg"`
	Data QueryResultData `json:"data"`
}

type QueryResultData struct {
	Qd string `json:"qd"`
	Fpdm string `json:"fpdm"`
	Fphm string `json:"fphm"`
	Kprq string `json:"kprq"`
	YzmSj string `json:"yzmSj"`
	Fpzt string `json:"fpzt"`
	Fxqy string `json:"fxqy"`
	Fplx string `json:"fplx"`
	Jqbm string `json:"jqbm"`
	Jym string `json:"jym"`
	GfName string `json:"gfName"`
	GfNsrsbh string `json:"gfNsrsbh"`
	GfAddressTel string `json:"gfAddressTel"`
	GfBankZh string `json:"gfBankZh"`
	JshjL string `json:"jshjL"`
	SfName string `json:"sfName"`
	SfNsrsbh string `json:"sfNsrsbh"`
	SfAddressTel string `json:"sfAddressTel"`
	SfBankZh string `json:"sfBankZh"`
	Bz string `json:"bz"`
	JshjU string `json:"jshjU"`
	ZpList []*ZpListItem `json:"zpList"`
}

type ZpListItem struct {
	MxName string `json:"mxName"`
	Ggxh string `json:"ggxh"`
	Unit string `json:"unit"`
	Num string `json:"num"`
	Price string `json:"price"`
	Je string `json:"je"`
	Sl string `json:"sl"`
	Se string `json:"se"`
}

type ResultItem struct {
	ResultId int `json:"resultid"`
	Ensured string `json:"ensured"`
	Sealed string `json:"sealed"`
	RespCode string `json:"respCode"`
	RespMsg string `json:"respMsg"`
	Qd string `json:"qd"`
	Fpdm string `json:"fpdm"`
	Fphm string `json:"fphm"`
	Kprq string `json:"kprq"`
	YzmSj string `json:"yzmSj"`
	Fpzt string `json:"fpzt"`
	Fxqy string `json:"fxqy"`
	Fplx string `json:"fplx"`
	Jqbm string `json:"jqbm"`
	Jym string `json:"jym"`
	GfName string `json:"gfName"`
	GfNsrsbh string `json:"gfNsrsbh"`
	GfAddressTel string `json:"gfAddressTel"`
	GfBankZh string `json:"gfBankZh"`
	JshjL string `json:"jshjL"`
	SfName string `json:"sfName"`
	SfNsrsbh string `json:"sfNsrsbh"`
	SfAddressTel string `json:"sfAddressTel"`
	SfBankZh string `json:"sfBankZh"`
	Bz string `json:"bz"`
	JshjU string `json:"jshjU"`
	MxName string `json:"mxName"`
	Ggxh string `json:"ggxh"`
	Unit string `json:"unit"`
	Price string `json:"price"`
	Je string `json:"je"`
	Sl string `json:"sl"`
	Se string `json:"se"`
	TotalJe float64 `json:"totalJe"`
	TotalSe float64 `json:"totalSe"`
	QueryTime string `json:"queryTime"`
	Num string `json:"num"`
}

func (item *ResultItem) DeepCopy() *ResultItem {
	copyItem := *item
	return &copyItem
}

type ResultItemWithUsername struct {
	ResultId int `json:"resultid"`
	Ensured string `json:"ensured"`
	Sealed string `json:"sealed"`
	RespCode string `json:"respCode"`
	RespMsg string `json:"respMsg"`
	Qd string `json:"qd"`
	Fpdm string `json:"fpdm"`
	Fphm string `json:"fphm"`
	Kprq string `json:"kprq"`
	YzmSj string `json:"yzmSj"`
	Fpzt string `json:"fpzt"`
	Fxqy string `json:"fxqy"`
	Fplx string `json:"fplx"`
	Jqbm string `json:"jqbm"`
	Jym string `json:"jym"`
	GfName string `json:"gfName"`
	GfNsrsbh string `json:"gfNsrsbh"`
	GfAddressTel string `json:"gfAddressTel"`
	GfBankZh string `json:"gfBankZh"`
	JshjL string `json:"jshjL"`
	SfName string `json:"sfName"`
	SfNsrsbh string `json:"sfNsrsbh"`
	SfAddressTel string `json:"sfAddressTel"`
	SfBankZh string `json:"sfBankZh"`
	Bz string `json:"bz"`
	JshjU string `json:"jshjU"`
	MxName string `json:"mxName"`
	Ggxh string `json:"ggxh"`
	Unit string `json:"unit"`
	Price string `json:"price"`
	Je string `json:"je"`
	Sl string `json:"sl"`
	Se string `json:"se"`
	TotalJe float64 `json:"totalJe"`
	TotalSe float64 `json:"totalSe"`
	QueryTime string `json:"queryTime"`
	Num string `json:"num"`
	Username string `json:"username"`
}

type ResultPara struct {
	UserName string `json:"username"`
	Filter string `json:"filter"`
	ResultId string `json:"resultid"`
	Operation string `json:"operation"`
	Rows int `json:"rows"`
	Page int `json:"page"`
}
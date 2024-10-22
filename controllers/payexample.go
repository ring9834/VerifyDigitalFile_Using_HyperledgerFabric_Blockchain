package controllers

import (
	"hzx/models"
)

func (c *MainController) PayQuery() {
	AccountID, _ := c.GetInt64("AccountID1")
	payment := models.GetPaymenRec(AccountID)
	c.Data["AccountID"] = payment.AccountID
	c.Data["PartnerID"] = payment.PartnerID
	c.Data["UserID"] = payment.UserID
	c.Data["CreateTime"] = payment.CreateTime
	c.Data["Amount"] = payment.Amount
	c.Data["OuterTradeNo"] = payment.OuterTradeNo
	c.Data["Remark"] = payment.Remark
	c.Data["Status"] = payment.Status
	c.Data["Msg"] = payment.Msg
	c.TplName = "query.html"
}
func (c *MainController) PayAdd() {
	var payment models.PaymentRecordStr
	c.ParseForm(&payment)
	pay := models.AddPaymenRec(payment)
	c.Data["AccountID"] = pay.AccountID
	c.Data["PartnerID"] = pay.PartnerID
	c.Data["UserID"] = pay.UserID
	c.Data["CreateTime"] = pay.CreateTime
	c.Data["Amount"] = pay.Amount
	c.Data["OuterTradeNo"] = pay.OuterTradeNo
	c.Data["Remark"] = pay.Remark
	c.TplName = "query.html"
}

func (c *MainController) GetPaysView() {
	c.TplName = "pays.html"
}

func (c *MainController) GetPays() {
	s := models.GetPayments()
	c.Data["json"] = s
	c.ServeJSON() //返回JSON对象
}

func (c *MainController) GetPagedPays() {
	pageIndex, _ := c.GetInt("pageIndex")
	pageSize, _ := c.GetInt("pageSize")
	rows, _, recordcount := models.GetPagedPayments(pageIndex, pageSize)

	rlt := make(map[string]interface{}, 1)
	rlt["total"] = recordcount //记录总条数,必须用total这个名称
	rlt["rows"] = rows         //当前页的数据记录,必须用rows这个名称
	c.Data["json"] = rlt
	c.ServeJSON() //返回JSON对象,供bootstrap绑定
}

package dbop

import (
	"receipts/model"
)

// 新增结果
func AddResult(username string, results []*model.ResultItem) error {

	stmt, err := db.Prepare("INSERT INTO result_" + username + "(ensured, sealed, respCode, respMsg, qd, fpdm, fphm, " +
		"kprq, yzmSj, fpzt, fxqy, fplx, jqbm, jym, gfName, gfNsrsbh, gfAddressTel, gfBankZh, jshjL, sfName, sfNsrsbh, " +
		"sfAddressTel, sfBankZh, bz, jshjU, mxName, ggxh, unit, price, je, sl, se, totalJe, totalSe, queryTime) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, item := range results {

		_, err = stmt.Exec(item.Ensured, item.Sealed, item.RespCode, item.RespMsg, item.Qd, item.Fpdm,
			item.Fphm, item.Kprq, item.YzmSj, item.Fpzt, item.Fxqy, item.Fplx, item.Jqbm, item.Jym,
			item.GfName, item.GfNsrsbh, item.GfAddressTel, item.GfBankZh, item.JshjL, item.SfName, item.SfNsrsbh,
			item.SfAddressTel, item.SfBankZh, item.Bz, item.JshjU, item.MxName, item.Ggxh, item.Unit, item.Price,
			item.Je, item.Sl, item.Se, item.TotalJe, item.TotalSe, item.QueryTime)

		if err != nil {
			return err
		}
	}

	return nil
}

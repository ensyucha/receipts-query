package dbop

import (
	"receipts/model"
)

// 新增结果
func AddResult(username string, result *model.ResultItem) error {

	stmt, err := db.Prepare("INSERT INTO result_" + username + "(ensured, sealed, respCode, respMsg, qd, fpdm, fphm, " +
		"kprq, yzmSj, fpzt, fxqy, fplx, jqbm, jym, gfName, gfNsrsbh, gfAddressTel, gfBankZh, jshjL, sfName, sfNsrsbh, " +
		"sfAddressTel, sfBankZh, bz, jshjU, zpListString) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(result.Ensured, result.Sealed, result.RespCode, result.RespMsg, result.Qd, result.Fpdm,
		result.Fphm, result.Kprq, result.YzmSj, result.Fpzt, result.Fxqy, result.Fplx, result.Jqbm, result.Jym,
		result.GfName, result.GfNsrsbh, result.GfAddressTel, result.GfBankZh, result.JshjL, result.SfName, result.SfNsrsbh,
		result.SfAddressTel, result.SfBankZh, result.Bz, result.JshjU, result.ZpListString)

	if err != nil {
		return err
	}

	return nil
}

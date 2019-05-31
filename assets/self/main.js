document.write('<script type="text/javascript" src="/assets/easyui/jquery.min.js"></script>');
document.write('<script type="text/javascript" src="/assets/easyui/jquery.easyui.min.js"></script>');
document.write('<script type="text/javascript" src="/assets/easyui/locale/easyui-lang-zh_CN.js"></script>');
document.write('<script type="text/javascript" src="/assets/other/jquery.cookie.js"></script>');
document.write('<script type="text/javascript" src="/assets/other/xlsx.core.min.js"></script>');
document.write('<script type="text/javascript" src="/assets/other/jquery-sortable.js"></script>');
document.write('<script type="text/javascript" src="/assets/mdui/js/mdui.js"></script>');

document.write('<link rel="icon" href="/assets/document/favicon.ico" />');
document.write('<link rel="stylesheet" href="/assets/easyui/themes/default/easyui.css">');
document.write('<link rel="stylesheet" href="/assets/easyui/themes/icon.css">');
document.write('<link rel="stylesheet" href="/assets/mdui/css/mdui.css">');
document.write('<link rel="stylesheet" href="/assets/self/main.css">');

/////////////////////////////////////////

function toResult() {
    window.location.href = "/result";
}

function toLogout() {
    window.location.href = "/logout";
}

function toQuery() {
    window.location.href = "/query";
}

function toSystem() {
    window.location.href = "/system";
}

function toTemp() {
    window.location.href = "/query/temp";
}

function toSealed() {
    window.location.href = "/result/sealed";
}

//////////////////////////////////////////

function showGuide(guideName) {
    window.open("/assets/document/"+guideName, "_blank");
}

function pagerFilter(data){
    if (typeof data.length == 'number' && typeof data.splice == 'function'){    // 判断数据是否是数组
        data = {
            total: data.length,
            rows: data
        }
    }
    let dg = $(this);
    let opts = dg.datagrid('options');
    let pager = dg.datagrid('getPager');
    pager.pagination({
        onSelectPage:function(pageNum, pageSize){
            opts.pageNumber = pageNum;
            opts.pageSize = pageSize;
            pager.pagination('refresh',{
                pageNumber:pageNum,
                pageSize:pageSize
            });
            dg.datagrid('loadData',data);
        }
    });
    if (!data.originalRows){
        data.originalRows = (data.rows);
    }
    let start = (opts.pageNumber-1)*parseInt(opts.pageSize);
    let end = start + parseInt(opts.pageSize);
    data.rows = (data.originalRows.slice(start, end));
    return data;
}

function outputExcel(jsonObjList) {

    let jsonObjListData = jsonObjList[0];

    jsonObjList.unshift({});

    let keyMap = []; //获取keys

    for (let key in jsonObjListData) {
        keyMap.push(key);
        jsonObjList[0][key] = key;
    }

    let savedData = [];//用来保存转换好的json

    jsonObjList.map((v, i) => keyMap.map((k, j) => Object.assign({}, {
        v: v[k],
        position: (j > 25 ? getCharCol(j) : String.fromCharCode(65 + j)) + (i + 1),
    }))).reduce((prev, next) => prev.concat(next)).forEach((v, i) => savedData[v.position] = {
        v: v.v
    });

    let outputPos = Object.keys(savedData); //设置区域,比如表格从A1到D10

    let tmpWB = {
        SheetNames: ['mySheet'], //保存的表标题
        Sheets: {
            'mySheet': Object.assign({},
                savedData, //内容
                {
                    '!ref': outputPos[0] + ':' + outputPos[outputPos.length - 1] //设置填充区域
                })
        }
    };

    let downloadObj = new Blob([stringToCharStream(XLSX.write(tmpWB,
        {bookType: "xlsm", bookSST: false, type: 'binary'}//这里的数据是用来定义导出的格式类型
    ))], {
        type: ""
    }); //创建二进制对象写入转换好的字节流

    document.getElementById("hf").href = URL.createObjectURL(downloadObj);

    document.getElementById("hf").click(); //模拟点击实现下载

    setTimeout(function() { //延时释放
        URL.revokeObjectURL(downloadObj); //用URL.revokeObjectURL()来释放这个object URL
    }, 100);
}

function outputSomeExcel(jsonObjList) {

    let jsonObjListData = jsonObjList[0];

    jsonObjList.unshift({});

    let keyMap = []; //获取keys

    for (let key in jsonObjListData) {
        keyMap.push(key);
        jsonObjList[0][key] = key;
    }

    let savedData = [];//用来保存转换好的json

    jsonObjList.map((v, i) => keyMap.map((k, j) => Object.assign({}, {
        v: v[k],
        position: (j > 25 ? getCharCol(j) : String.fromCharCode(65 + j)) + (i + 1),
    }))).reduce((prev, next) => prev.concat(next)).forEach((v, i) => savedData[v.position] = {
        v: v.v
    });

    let outputPos = Object.keys(savedData); //设置区域,比如表格从A1到D10

    let tmpWB = {
        SheetNames: ['mySheet'], //保存的表标题
        Sheets: {
            'mySheet': Object.assign({},
                savedData, //内容
                {
                    '!ref': outputPos[0] + ':' + outputPos[outputPos.length - 1] //设置填充区域
                })
        }
    };

    let downloadObj = new Blob([stringToCharStream(XLSX.write(tmpWB,
        {bookType: "xlsm", bookSST: false, type: 'binary'}//这里的数据是用来定义导出的格式类型
    ))], {
        type: ""
    }); //创建二进制对象写入转换好的字节流

    document.getElementById("hfsome").href = URL.createObjectURL(downloadObj);

    document.getElementById("hfsome").click(); //模拟点击实现下载

    setTimeout(function() { //延时释放
        URL.revokeObjectURL(downloadObj); //用URL.revokeObjectURL()来释放这个object URL
    }, 100);
}

function stringToCharStream(str) {
    let buf = new ArrayBuffer(str.length);

    let view = new Uint8Array(buf);

    for (let i = 0; i !== str.length; ++i) {
        view[i] = str.charCodeAt(i) & 0xFF;
    }

    return buf;
}

function getCharCol(n) {

    let s = '';
    let m = 0;

    while (n > 0) {
        m = n % 26 + 1;
        s = String.fromCharCode(m + 64) + s;
        n = (n - m) / 26;
    }

    return s;
}

function makeAllResultOutput(item) {

    let zpList = item['zpListString'].split("||");

    return {
        "发票类型": item['fplx'],
        "发票代码": item['fpdm'],
        "发票号码": item['fphm'],
        "开票日期": item['kprq'],
        "购方名称": item['gfName'],
        "购方税号": item['gfNsrsbh'],
        "购方地址电话": item['gfAddressTel'],
        "购方开户行账号": item['gfBankZh'],
        "销售方名称": item['sfName'],
        "销方纳税人识别号": item['sfNsrsbh'],
        "销方地址及电话": item['sfAddressTel'],
        "销方银行及账号": item['sfBankZh'],
        "不含税金额": zpList[5],
        "税额": zpList[7],
        "发票明细": zpList[0],
        "规格型号": zpList[1],
        "计量单位": zpList[4],
        "数量": zpList[3],
        "单价": zpList[2],
        "税率": zpList[6],
        "发票备注": item['bz'],
    }
}

function makeSomeResultOutput(item) {

    let zpList = item['zpListString'].split("||");

    return {

        "发票代码": item['fpdm'],
        "发票号码": item['fphm'],
        "开票日期": item['kprq'],
        "销方纳税人识别号": item['sfNsrsbh'],
        "销售方名称": item['sfName'],
        "金额": zpList[5],
        "税额": zpList[7],
        "认证方式": "???",
        "确认/认证日期": item['yzmSj'],
        "发票类型": item['fplx'],
        "发票状态": item['fpzt'],
    }
}

/////////////////////////////////////////

function login() {

    let input_username = $("#input-username").val();
    let input_password = $("#input-password").val();

    let errInfo = "";

    if(input_username.length <= 0) {
        errInfo += "1.请输入账号<br />";
    }

    if(input_password.length <= 0) {
        errInfo += "2.请输入密码<br />";
    }

    if (errInfo.length === 0) {

        $.ajax({ // 发起 ajax 请求，进行密码判断
            type: "POST",
            timeout: 5000,
            url: "/login",
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            data: JSON.stringify({
                "username": input_username,
                "password": input_password
            }),
            success: function (message) {

                mdui.snackbar({
                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                    position: "top",
                    timeout: 100,
                    onClose: function(){
                        if (input_username === "admin") {
                            toSystem();
                        } else {
                            toQuery();
                        }
                    }
                });
            },
            error: function(message) {
                mdui.snackbar({
                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                    position: "top",
                    timeout: 3000,
                });
            }
        });
    } else {
        mdui.snackbar({
            message: errInfo, // 显示错误信息给用户
            position: "top",
            timeout: 2000,
        });
    }
}

/////////////////////////////////////////

function extractExcelData(e) {

    let files = e.target.files;
    let fileReader = new FileReader();

    let filename = excelFileSelector.val();
    let suffix = (filename.slice(filename.length - 4));

    if (filename.length < 4 || suffix !== "xlsm") {
        mdui.snackbar({
            message: "文件必须是Excel表格，且后缀名为 .xlsm ",
            position: "top",
            timeout: 2000,
        });
        return; // 结束函数
    }

    fileReader.onload = function(ev) {

        let workbook; // 存储 excel 表格对象
        let excelData = []; // 存储获取到的数据

        try {
            let data = ev.target.result; // 获取上传文件的数据
            workbook = XLSX.read(data, { // 以二进制流方式读取得到整份 excel 表格对象
                type: 'binary'
            });
        } catch (e) {
            mdui.snackbar({
                message: "文件类型不正确",
                position: "top",
                timeout: 1000,
            });
            return; // 结束函数
        }

        // 遍历每张表读取数据
        for (let sheet in workbook.Sheets) {
            if (workbook.Sheets.hasOwnProperty(sheet)) {
                excelData = excelData.concat(XLSX.utils.sheet_to_json(workbook.Sheets[sheet])); // 获取数据
                break; // 只读取第一张表的内容
            }
        }

        processExcelReceipt(excelData); // 处理上传的Excel表格数据

        dgSelector.datagrid('loadData', queryrows);
    };

    fileReader.readAsBinaryString(files[0]); // 以二进制方式打开文件

    $("#excel-file").val(""); // 清空已上传文件，避免无法上传同一文件多次
}

function processExcelReceipt(excelData) {

    for (let i=1; i<excelData.length; i++) { // 遍历 excelData

        let fpdm = excelData[i]["发票代码(必填)"];
        let fphm = excelData[i]["发票号码(必填)"];
        let kprq = excelData[i]["开票日期(必填)"];
        let je = excelData[i]["不含税金额(必填)"];
        let jym = excelData[i]["校验码后6位(选填)"];

        if (!je) { je = ""; }
        if (!jym) { jym = ""; }

        let obj = {
            "fpdm": fpdm.toString(),
            "fphm": fphm.toString(),
            "kprq": kprq.toString(),
            "je"  : je.toString(),
            "jym": jym.toString()
        };

        queryrows.push(obj);
    }
}

function addByExcel() {
    $("#excel-file").click();
}

function query(usage) {
    usage = parseInt(usage);

    if (usage < queryrows.length) {
        mdui.snackbar({
            message: "额度不足，请增加额度或减少查询量", // 显示信息给用户
            position: "top",
            timeout: 3000,
        });
    } else {

        mdui.confirm('剩余额度：' + usage + '<br />查询发票数量：' + queryrows.length + '<br />', "", function() {

            mdui.snackbar({
                message: '查询中，请稍等', // 显示信息给用户
                position: "top",
                timeout: 500,
            });

            $.ajax({ // 发起 ajax 请求，进行密码判断
                type: "POST",
                timeout: 20000,
                url: "/query",
                contentType: "application/json; charset=utf-8",
                dataType: "json",
                data: JSON.stringify({
                    "queryArray": queryrows
                }),
                success: function (message) {

                    mdui.snackbar({
                        message: message['status'] + " : " + message['message'], // 显示信息给用户
                        position: "top",
                        timeout: 1000,
                        onClose: function () {
                            toTemp();
                        }
                    });
                },
                error: function (message) {
                    toTemp();
                }
            });

        },function () {},{confirmText:"确认查询",cancelText:"取消",modal:true,closeOnEsc:false});

    }
}

function clearQueryData() {
    queryrows = [];
    $("#dg").datagrid('loadData', queryrows);
}

/////////////////////////////////////////

function initTemp(tempDataList) {
    let rows = [];

    for (let i = 0; i < tempDataList.length; i++) {
        if (tempDataList[i]['data']['zpList'] != null) {
            tempDataList[i]['data']['bz'] = tempDataList[i]['data']['zpList'][0]['je'];
        }

        let querySuccess = "<span class='my-failed'>查询失败，不扣额度</span>";
        let queryRespMsg = "<span class='my-failed'>" + tempDataList[i]['respMsg'] + "</span>";

        let respCode = tempDataList[i]['respCode'];
        if (respCode === "2210" || respCode === "2213" || respCode === "2215" || respCode === "2206") {
            querySuccess = "<span class='my-success'>查询成功，扣除额度</span>";
            queryRespMsg = "<span class='my-success'>" + tempDataList[i]['respMsg'] + "</span>";
        }

        rows.push({
            status: querySuccess,
            respMsg: queryRespMsg,
            fpdm: tempDataList[i]['data']['fpdm'],
            fphm: tempDataList[i]['data']['fphm'],
            kprq: tempDataList[i]['data']['kprq'],
            bz:   tempDataList[i]['data']['bz'],
            jym:  tempDataList[i]['data']['jym']
        })
    }

    return rows;
}

/////////////////////////////////////////

function updatePassword(curPassword) {

    mdui.dialog({
        title: '<span class="dialog-title-color">修改管理密码</span>',
        content: `当前密码：<b>` + curPassword + `</b>` +
            `<div class="mdui-textfield mdui-textfield-floating-label">
                        <label class="mdui-textfield-label">请输入新密码</label>
                        <input id="input-password" class="mdui-textfield-input" type="password" />
                    </div>`,
        buttons: [
            {text: '关闭'},
            {
                text: '更新',
                onClick: function (inst) {

                    let password = $("#input-password").val(); // 获取密码

                    let regExp = new RegExp(/^[0-9a-zA-Z]{8,}$/);

                    if (!regExp.test(password)) {
                        mdui.snackbar({
                            message: "密码只能含有数字和英文字母，且长度要大于等于8", // 显示错误信息给用户
                            position: "top",
                            timeout: 5000,
                        });
                    } else {

                        mdui.snackbar({ // 显示 更新中 snackbar
                            message: "更新中...",
                            position: "top",
                            timeout: 500,
                        });

                        $.ajax({ // 发起 ajax 请求，进行密码判断
                            type: "POST",
                            timeout: 5000,
                            url: "/system/password",
                            contentType: "application/json; charset=utf-8",
                            dataType: "json",
                            data: JSON.stringify({"password": password}),
                            success: function (message) {
                                inst.close(); // 关闭表单窗口
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 500,
                                    onClose: function(){
                                        toLogout();
                                    }
                                });
                            },
                            error: function(message) {
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 3000,
                                });
                            }
                        });
                    }
                },
                close: false, // 禁止直接点击按钮关闭窗口
            }
        ],
        modal: true,
        closeOnEsc: false,
    });
}

function updateUnusedUsage(curUnusedUsage, apicode) {

    if (apicode !== "empty") {
        mdui.dialog({
            title: '<span class="dialog-title-color">更新未分配余额</span>',
            content: `当前未分配额度：<b>` + curUnusedUsage + `</b>` +
                `<div class="mdui-textfield mdui-textfield-floating-label">
                        <label class="mdui-textfield-label">请输入新的未分配额度</label>
                        <input id="input-unusedusage" class="mdui-textfield-input" autocomplete="off" />
                    </div>`,
            buttons: [
                {text: '关闭'},
                {
                    text: '更新',
                    onClick: function (inst) {

                        let unusedUsage = parseInt($("#input-unusedusage").val());
                        let unusedUsageReg = /^[1-9][0-9]*$/;
                        let regExp = new RegExp(unusedUsageReg);

                        if (!regExp.test(unusedUsage + "") && unusedUsage !== 0) {
                            mdui.snackbar({
                                message: "请输入新的余额<br />只能是数字，且不能以0开头<br />", // 显示错误信息给用户
                                position: "top",
                                timeout: 5000,
                            });
                        } else {

                            mdui.snackbar({ // 显示 更新中 snackbar
                                message: "更新中...",
                                position: "top",
                                timeout: 500,
                            });

                            $.ajax({ // 发起 ajax 请求，进行密码判断
                                type: "POST",
                                timeout: 5000,
                                url: "/system/unusedusage",
                                contentType: "application/json; charset=utf-8",
                                dataType: "json",
                                data: JSON.stringify({"unusedusage": unusedUsage}),
                                success: function (message) {
                                    inst.close(); // 关闭表单窗口
                                    mdui.snackbar({
                                        message: message['status'] + " : " + message['message'], // 显示信息给用户
                                        position: "top",
                                        timeout: 500,
                                        onClose: function () {
                                            window.location.reload();
                                        }
                                    });
                                },
                                error: function (message) {
                                    mdui.snackbar({
                                        message: message['status'] + " : " + message['message'], // 显示信息给用户
                                        position: "top",
                                        timeout: 3000,
                                    });
                                }
                            });
                        }
                    },
                    close: false, // 禁止直接点击按钮关闭窗口
                }
            ],
            modal: true,
            closeOnEsc: false,
        });
    } else {
        mdui.snackbar({
            message: "设置ApiCode后才能设置未分配余额", // 显示信息给用户
            position: "top",
            timeout: 3000,
        });
    }
}

function updateApiCode(curApiCode) {

    mdui.dialog({
        title: '<span class="dialog-title-color">更新ApiCode</span>',
        content: `当前ApiCode：<b>` + curApiCode + `</b>` +
            `<div class="mdui-textfield mdui-textfield-floating-label">
                        <label class="mdui-textfield-label">请输入新的apicode</label>
                        <input id="input-apicode" class="mdui-textfield-input" autocomplete="off" />
                    </div>`,
        buttons: [
            {text: '关闭'},
            {
                text: '更新',
                onClick: function (inst) {

                    let apicode = $("#input-apicode").val(); // 获取密码

                    let regExp = new RegExp(/^[a-zA-Z0-9]+$/);

                    if (!regExp.test(apicode)) {
                        mdui.snackbar({
                            message: "请输入正确的apicode", // 显示错误信息给用户
                            position: "top",
                            timeout: 3000,
                        });
                    } else {

                        mdui.snackbar({ // 显示 更新中 snackbar
                            message: "更新中...",
                            position: "top",
                            timeout: 500,
                        });

                        $.ajax({ // 发起 ajax 请求，进行密码判断
                            type: "POST",
                            timeout: 5000,
                            url: "/system/apicode",
                            contentType: "application/json; charset=utf-8",
                            dataType: "json",
                            data: JSON.stringify({"apicode": apicode}),
                            success: function (message) {
                                inst.close(); // 关闭表单窗口
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 500,
                                    onClose: function(){
                                        window.location.reload();
                                    }
                                });
                            },
                            error: function(message) {
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 3000,
                                });
                            }
                        });

                    }
                },
                close: false, // 禁止直接点击按钮关闭窗口
            }
        ],
        modal: true,
        closeOnEsc: false,
    });
}

function addUser() {
    mdui.dialog({
        title: '<span class="dialog-title-color">新增用户</span>',
        content: `
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">账号</label>
                    <input id="input-username" class="mdui-textfield-input" autocomplete="off" />
                </div>
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">昵称</label>
                    <input id="input-nickname" class="mdui-textfield-input" autocomplete="off" />
                </div>
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">密码</label>
                    <input id="input-password" class="mdui-textfield-input" autocomplete="off" />
                </div>
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">分配余额</label>
                    <input id="input-usage" class="mdui-textfield-input" autocomplete="off" />
                </div>`,
        buttons: [
            { text: '取消' },
            {
                text: '新增',
                onClick: function(inst){
                    let result = checkAddUser();

                    if (result.errInfo.length === 0) {

                        mdui.snackbar({ // 显示 更新中 snackbar
                            message: "更新中...",
                            position: "top",
                            timeout: 500,
                        });

                        $.ajax({ // 发起 ajax 请求，进行密码判断
                            type: "POST",
                            timeout: 5000,
                            url: "/system/user",
                            contentType: "application/json; charset=utf-8",
                            dataType: "json",
                            data: JSON.stringify({
                                "username": result.username,
                                "nickname": result.nickname,
                                "password": result.password,
                                "usage": result.usage
                            }),
                            success: function (message) {
                                inst.close(); // 关闭表单窗口
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 500,
                                    onClose: function(){
                                        window.location.reload()
                                    }
                                });
                            },
                            error: function(message) {
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 3000,
                                });
                            }
                        });

                    } else { // 如果单行发票填写格式不正确
                        mdui.snackbar({
                            message: result.errInfo, // 显示错误信息给用户
                            position: "top",
                            timeout: 7000,
                        });
                    }
                },
                close: false, // 禁止直接点击按钮关闭窗口
            }
        ],
        modal: true, // 禁止点击空白区域关闭窗口
        closeOnEsc: false, // 禁止通过ESC关闭窗口
    });
}

function checkAddUser() {

    let input_username = $("#input-username").val();
    let input_nickname = $("#input-nickname").val();
    let input_password = $("#input-password").val();
    let input_usage = $("#input-usage").val();

    let errInfo = "";

    let regExp1 = new RegExp(/^[a-zA-Z][a-zA-Z0-9]{5,18}$/);
    if(!regExp1.test(input_username)) {
        errInfo += "1.账号长度要大于5且小于20，账号只能含有英文字母或数字，且只能以英文字母开头<br />";
    }

    if(input_nickname.length < 6 || input_nickname.length >= 20) {
        errInfo += "2.用户昵称长度要大于5且小于20<br />";
    }

    let regExp2 = new RegExp(/^[0-9a-zA-Z]{6,19}$/);
    if(!regExp2.test(input_password)) {
        errInfo += "3.密码只能含有数字和英文字母，且长度要大于5且小于20<br />";
    }

    let regExp3 = new RegExp(/^[1-9][0-9]*$/);
    if (input_usage !== "0") {
        if (!regExp3.test(input_usage)) {
            errInfo += "4.额度分配必须输入整数<br />";
        }
    }

    return { // 返回错误信息和单条发票信息
        "errInfo": errInfo,
        "username": input_username,
        "nickname": input_nickname,
        "password": input_password,
        "usage": parseInt(input_usage)
    };
}

function listAllUser() {

    $.ajax({ // 发起 ajax 请求，进行密码判断
        type: "GET",
        timeout: 5000,
        url: "/system/user",
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (message) {
            let userList = message['message'];

            for (let i=0; i<userList.length; i++) {
                let userObj = {
                    "username": userList[i].username,
                    "nickname": userList[i].nickname,
                    "password": userList[i].password,
                    "usage": userList[i].usage
                };
                listOneUser(userObj);
            }
        },
        error: function(message) {
            mdui.snackbar({
                message: message['status'] + " : " + message['message'], // 显示信息给用户
                position: "top",
                timeout: 3000,
            });
        }
    });
}

function listOneUser(userObj) {

    let userObjStr = JSON.stringify(userObj);
    userObjStr = userObjStr.replace(/"/g, '@@');

    let t = `<tr>
                <td class="system-td-center">` + userObj.username + `</td>
                <td class="system-td-center">` + userObj.nickname + `</td>
                <td class="system-td-center">` + userObj.password + `</td>
                <td class="system-td-center">` + userObj.usage + `</td>
                <td class="system-td-center">
                    <button class="mdui-btn mdui-color-red" onclick="editUser('` + userObjStr + `')">编辑</button>
                    <button class="mdui-btn mdui-color-red" onclick="removeUser('` + userObjStr + `')">删除</button>
                </td>
            </tr>`;

    $("#main-body").append(t)
}

function removeUser(userObjStr) {

    userObjStr = userObjStr.replace(/@@/g, '"');
    let userObj = JSON.parse(userObjStr);

    mdui.confirm('确定要删除用户吗？<br />一旦删除其数据将被清空！<br />请谨慎操作！', "",function() {

        $.ajax({ // 发起 ajax 请求，进行密码判断
            type: "DELETE",
            timeout: 5000,
            url: "/system/user",
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            data: JSON.stringify({
                "username": userObj.username
            }),
            success: function (message) {
                mdui.snackbar({
                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                    position: "top",
                    timeout: 500,
                    onClose: function(){
                        window.location.reload()
                    }
                });
            },
            error: function(message) {
                mdui.snackbar({
                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                    position: "top",
                    timeout: 3000,
                });
            }
        });

    },function () {},{confirmText:"确认删除用户",cancelText:"取消",modal:true,closeOnEsc:false});
}

function editUser(userObjStr) {
    userObjStr = userObjStr.replace(/@@/g, '"');
    let userObj = JSON.parse(userObjStr);

    mdui.dialog({
        title: '<span class="dialog-title-color">编辑用户</span>',
        content: `
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">昵称</label>
                    <input id="edit-nickname" class="mdui-textfield-input" autocomplete="off" />
                </div>
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">密码</label>
                    <input id="edit-password" class="mdui-textfield-input" autocomplete="off" />
                </div>
                <div class="mdui-textfield mdui-textfield-floating-label">
                    <label class="mdui-textfield-label">分配余额</label>
                    <input id="edit-usage" class="mdui-textfield-input" autocomplete="off" />
                </div>`,
        buttons: [
            { text: '取消' },
            {
                text: '更新',
                onClick: function(inst){

                    let errInfo = "";
                    let edit_nickname = $("#edit-nickname").val();
                    let edit_password = $("#edit-password").val();
                    let edit_usage = $("#edit-usage").val();

                    if(edit_nickname.length < 6 || edit_nickname.length >= 20) {
                        errInfo += "1.用户昵称长度要大于5且小于20<br />";
                    }

                    let regExp1 = new RegExp(/^[0-9a-zA-Z]{6,19}$/);
                    if(!regExp1.test(edit_password)) {
                        errInfo += "2.密码只能含有数字和英文字母，且长度要大于5且小于20<br />";
                    }

                    let regExp2 = new RegExp(/^[1-9][0-9]*$/);
                    if (edit_usage !== "0") {
                        if (!regExp2.test(edit_usage)) {
                            errInfo += "3.额度分配必须输入整数<br />";
                        }
                    }

                    if (errInfo.length === 0) {

                        mdui.snackbar({ // 显示 更新中 snackbar
                            message: "更新中...",
                            position: "top",
                            timeout: 500,
                        });

                        $.ajax({ // 发起 ajax 请求，进行密码判断
                            type: "PUT",
                            timeout: 5000,
                            url: "/system/user",
                            contentType: "application/json; charset=utf-8",
                            dataType: "json",
                            data: JSON.stringify({
                                "username": userObj.username,
                                "nickname": edit_nickname,
                                "password": edit_password,
                                "usage": parseInt(edit_usage)
                            }),
                            success: function (message) {
                                inst.close(); // 关闭表单窗口
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 500,
                                    onClose: function(){
                                        window.location.reload()
                                    }
                                });
                            },
                            error: function(message) {
                                mdui.snackbar({
                                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                                    position: "top",
                                    timeout: 3000,
                                });
                            }
                        });

                    } else { // 如果单行发票填写格式不正确
                        mdui.snackbar({
                            message: errInfo, // 显示错误信息给用户
                            position: "top",
                            timeout: 7000,
                        });
                    }
                },
                close: false, // 禁止直接点击按钮关闭窗口
            }
        ],
        modal: true, // 禁止点击空白区域关闭窗口
        closeOnEsc: false, // 禁止通过ESC关闭窗口
        onOpened: function() {
            $("#edit-nickname").val(userObj.nickname);
            $("#edit-password").val(userObj.password);
            $("#edit-usage").val(userObj.usage);

            mdui.mutation();
        }
    });


}

/////////////////////////////////////////

function editShowHideCol() {
    mdui.dialog({
        title: '<span class="dialog-title-color">显示/隐藏列</span>',
        content: `<div class="result-dialog-content">打钩显示，不打勾隐藏</div><br />
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-ensured" value="ensured"/><i class="mdui-checkbox-icon"></i>确认状态</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-respMsg" value="respMsg"/><i class="mdui-checkbox-icon"></i>验证类型</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-fpzt" value="fpzt"/><i class="mdui-checkbox-icon"></i>发票状态</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-fplx" value="fplx"/><i class="mdui-checkbox-icon"></i>发票类型</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-jshjL" value="jshjL"/><i class="mdui-checkbox-icon"></i>价税合计</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-jshjU" value="jshjU"/><i class="mdui-checkbox-icon"></i>价税合计(大写)</label>
            <br /><br />
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-fpdm" value="fpdm"/><i class="mdui-checkbox-icon"></i>发票代码</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-fphm" value="fphm"/><i class="mdui-checkbox-icon"></i>发票号码</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-kprq" value="kprq"/><i class="mdui-checkbox-icon"></i>开票日期</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-yzmSj" value="yzmSj"/><i class="mdui-checkbox-icon"></i>验证时间</label>
            <br /><br />
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-jym" value="jym"/><i class="mdui-checkbox-icon"></i>检验码</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-qd" value="qd"/><i class="mdui-checkbox-icon"></i>有无清单</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-jqbm" value="jqbm"/><i class="mdui-checkbox-icon"></i>机器编码</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-zpListString" value="zpListString"/><i class="mdui-checkbox-icon"></i>商品列表</label>
            <br /><br />
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-gfName" value="gfName"/><i class="mdui-checkbox-icon"></i>购方名称</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-gfNsrsbh" value="gfNsrsbh"/><i class="mdui-checkbox-icon"></i>购方识别号</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-gfAddressTel" value="gfAddressTel"/><i class="mdui-checkbox-icon"></i>购方联系地址</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-gfBankZh" value="gfBankZh"/><i class="mdui-checkbox-icon"></i>购方开户行</label>
            <br /><br />
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-sfName" value="sfName"/><i class="mdui-checkbox-icon"></i>销售方名称</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-sfNsrsbh" value="sfNsrsbh"/><i class="mdui-checkbox-icon"></i>销售方识别号</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-sfAddressTel" value="sfAddressTel"/><i class="mdui-checkbox-icon"></i>销售方联系方式</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-sfBankZh" value="sfBankZh"/><i class="mdui-checkbox-icon"></i>销售方开户行</label>
            <br /><br />
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-fxqy" value="fxqy"/><i class="mdui-checkbox-icon"></i>风险企业验证</label>
            <label class="mdui-checkbox result-checkbox-label"><input class="cb" type="checkbox" id="checkbox-bz" value="bz"/><i class="mdui-checkbox-icon"></i>备注</label>
            `,
        buttons: [
            { text: '取消' },
            {
                text: '更新',
                onClick: function(inst){
                    showHideColPara['show'] = [];
                    showHideColPara['hide'] = [];

                    $.each($(".cb"), function(){
                        if ($(this).prop('checked')) {
                            showHideColPara['show'].push($(this).val())
                        } else {
                            showHideColPara['hide'].push($(this).val())
                        }
                    });

                    $.cookie('showHidePara', JSON.stringify(showHideColPara), { expires: 100 * 365 });

                    showHideCol();

                    dgSelector.datagrid("resize");

                    inst.close();
                },
                close: false, // 禁止直接点击按钮关闭窗口
            }
        ],
        modal: true, // 禁止点击空白区域关闭窗口
        closeOnEsc: false, // 禁止通过ESC关闭窗口
        onOpened: function() {
            let showPara = showHideColPara['show'];
            for (let i=0; i < showPara.length; i++) {
                $("#checkbox-"+showPara[i]).attr("checked",true);
            }
        }
    });
}

function showHideCol() {

    for (let i=0; i<showHideColPara['show'].length; i++) {
        dgSelector.datagrid('showColumn', showHideColPara['show'][i]);
    }
    for (let i=0; i<showHideColPara['hide'].length; i++) {
        dgSelector.datagrid('hideColumn', showHideColPara['hide'][i]);
    }

    dgSelector.datagrid("resize");
}

function buildCol(colInfo) {
    let col = [];

    for (let i=0; i<colInfo.length; i++) {
        let dgField = {
            field: colInfo[i][0],
            title: colInfo[i][1],
            hidden: !colInfo[i][2],
            width: colInfo[i][3],
            align: colInfo[i][4],
            sortable:true,
            sorter:function(a,b) { return (a < b ? 1 : -1); }
        };
        col.push(dgField)
    }

    return col;
}

function getSelectedResultID() {
    let ids = [];
    let src = dgSelector.datagrid("getSelections");

    for (let i=0; i<src.length; i++) {
        ids.push(src[i].resultid);
    }

    return ids;
}

function makeResultPara(username, filter,idsString,operation) {

    return "username=" + username + "&filter=" + filter + "&resultid=" + idsString +
    "&operation=" + operation + "&page=1&rows=20";
}

function removeResult(username) {

    let idsString = getSelectedResultID().join("-");

    mdui.confirm('确定要删除吗？<br />一旦删除将无法恢复！<br />请谨慎操作！', "",function() {

        $.ajax({ // 发起 ajax 请求，进行密码判断
            type: "POST",
            timeout: 5000,
            url: "/result/data",
            data: makeResultPara(username,"", idsString, "removedata"),
            success: function (message) {
                mdui.snackbar({
                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                    position: "top",
                    timeout: 100,
                    onClose: function() {
                        $('#dg').datagrid('reload');
                    }
                });
            },
            error: function(message) {
                mdui.snackbar({
                    message: message['status'] + " : " + message['message'], // 显示信息给用户
                    position: "top",
                    timeout: 3000,
                });
            }
        });

    },function () {},{confirmText:"确认删除",cancelText:"取消",modal:true,closeOnEsc:false});
}

function outputAllResult(username, filter) {

    let filterString = "";

    if (filter !== "sealed") {
        filterString = "sealed%3D'0'"
    } else {
        filterString = "sealed%3D'1'"
    }

    $.ajax({ // 发起 ajax 请求，进行密码判断
        type: "POST",
        timeout: 5000,
        url: "/result/data",
        data: makeResultPara(username, filterString, "", "getdata"),
        success: function (message) {
            mdui.snackbar({
                message: "获取数据成功，正在导出", // 显示信息给用户
                position: "top",
                timeout: 500,
                onClose: function() {
                    let jsonObjList = [];
                    let rows = message.rows;

                    for (let i=0; i<rows.length; i++) {
                        let item = rows[i];

                        jsonObjList.push(makeAllResultOutput(item))
                    }

                    outputExcel(jsonObjList);
                }
            });
        },
        error: function(message) {
            mdui.snackbar({
                message: message['status'] + " : " + message['message'], // 显示信息给用户
                position: "top",
                timeout: 3000,
            });
        }
    });
}

function outputSomeResult(username, filter) {

    let filterString = "";

    if (filter !== "sealed") {
        filterString = "sealed%3D'0'"
    } else {
        filterString = "sealed%3D'1'"
    }

    $.ajax({ // 发起 ajax 请求，进行密码判断
        type: "POST",
        timeout: 5000,
        url: "/result/data",
        data: makeResultPara(username, filterString, "", "getdata"),
        success: function (message) {
            mdui.snackbar({
                message: "获取数据成功，正在导出", // 显示信息给用户
                position: "top",
                timeout: 500,
                onClose: function() {
                    let jsonObjList = [];
                    let rows = message.rows;

                    for (let i=0; i<rows.length; i++) {
                        let item = rows[i];

                        jsonObjList.push(makeSomeResultOutput(item))
                    }

                    outputSomeExcel(jsonObjList);
                }
            });
        },
        error: function(message) {
            mdui.snackbar({
                message: message['status'] + " : " + message['message'], // 显示信息给用户
                position: "top",
                timeout: 3000,
            });
        }
    });
}

function updateResult(username, operation) {

    let idsString = getSelectedResultID().join("-");

    $.ajax({ // 发起 ajax 请求，进行密码判断
        type: "POST",
        timeout: 5000,
        url: "/result/data",
        data: makeResultPara(username,"", idsString, operation),
        success: function (message) {
            mdui.snackbar({
                message: message['status'] + " : " + message['message'], // 显示信息给用户
                position: "top",
                timeout: 100,
                onClose: function() {
                    if (operation === "sealed" || operation === "unsealed") {
                        window.location.reload();
                    } else {
                        let dgSelector = $('#dg');
                        dgSelector.datagrid('reload');
                        dgSelector.datagrid('unselectAll');
                    }
                }
            });
        },
        error: function(message) {
            mdui.snackbar({
                message: message['status'] + " : " + message['message'], // 显示信息给用户
                position: "top",
                timeout: 3000,
            });
        }
    });
}

function filterResult() {
    alert("待讨论细节");
}

////////////////////////////////////////////////////////////////////

let basicColInfo = [
    ['ensured','确认状态',true,'80px','center'],
    ['respMsg','发票状态',true,'140px','center'],
    ['fpzt','发票状态',true,'85px','center'],
    ['fplx','发票类型',true,'85px','center'],
    ['fpdm','发票代码',true,'95px','center'],
    ['fphm','发票号码',true,'85px','center'],
    ['kprq','开票日期',true,'90px','center'],
    ['yzmSj','验证时间',true,'145px','center'],
    ['jqbm','机器编码',true,'110px','center'],
    ['gfName','购方名称',true,'','center'],
    ['gfNsrsbh','购方识别号',true,'170px','center'],
    ['gfAddressTel','购方联系地址',true,'','center'],
    ['gfBankZh','购方开户行',true,'','center'],
    ['fxqy','风险企业验证',true,'92px','center'],
    ['sfName','销售方名称',true,'','center'],
    ['sfNsrsbh','销售方识别号',true,'170px','center'],
    ['sfAddressTel','销售方联系地址',true,'','center'],
    ['sfBankZh','销售方开户行',true,'','center'],
    ['jshjL','价税合计',true,'145px','center'],
    ['jshjU','价税合计(大写)',true,'','center'],
    ['qd','清单',true,'80px','center'],
    ['bz','备注',true,'','center'],
    ['zpListString','商品列表',true,'','center'],
    ['jym','校验码',true,'180px','center']
];
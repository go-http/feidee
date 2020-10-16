# feidee [![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/go-http/feidee/v2)](https://pkg.go.dev/mod/github.com/go-http/feidee/v2)
随手记API

# Usage

## 创建并初始化客户端
```go
	//创建客户端并登录
	client := feidee.New("username", "password")
	if err != nil {
		return fmt.Errorf("登录失败:%s", err)
	}

	//（可选）获取账本列表
	err = client.SyncAccountBookList()
	if err != nil {
		return fmt.Errorf("刷新账本列表失败: %s", err)
	}

	//（可选）切换账本（此处需要上一步同步账本列表才能通过账本名查到账本ID，否则会报错找不到账本）
	err = client.SwitchBook("默认账本")
	if err != nil {
		return fmt.Errorf("切换账本失败: %s", err)
	}

	//同步基础数据（科目、成员、商家、项目、账户等）
	err = client.SyncMetaInfo()
	if err != nil {
		return fmt.Errorf("同步账户基础信息失败:%s", err)
	}
```

## 打印上个月对账信息
```go
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	start := t.AddDate(0, -1, 0)
	end := start.AddDate(0, 1, -1)

	for id, account := range FdCLient.AccountMap {
		list, err := FdCLient.CompareReport(id, start, end)
		if err != nil {
			fmt.Printf("\n\n账户%s错误:%s", account.Name, err)
			continue
		}
		fmt.Println("\n\n账户:", account.Name, "x", len(list))
		for _, info := range list {
			fmt.Printf("\t%+v\n", info)
		}
	}
```

## 查询并打印按月收支
```go
	infoMap, err := FdCLient.MonthIncomeAndPayoutMap(2017, 2018)
	if err != nil {
		fmt.Printf("查询按月收支信息失败: %s", err)
		return
	}
	for key, info := range infoMap {
		fmt.Printf("%d: %+v\n", key, info)
	}
```

## 查询并打印上月流水
```go
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	start := t.AddDate(0, -1, 0)
	end := start.AddDate(0, 1, -1)

	info, err := FdCLient.TallyList(start, end, nil)
	if err != nil {
		fmt.Printf("查询失败:%s", err)
		return
	}

	fmt.Printf("%s - %s\n", info.BeginDate, info.EndDate)
	fmt.Printf("汇总：%+v\n", info.IncomeAndPayout)
	for i, group := range info.Groups {
		fmt.Printf("%+v\n", group.IncomeAndPayout)
		for j, t := range group.List {
			fmt.Printf("%d-%d %4s(%2d) %8s => %8s %10.2f\n", i, j, t.TranName, t.TranType, t.BuyerAcount, t.SellerAcount, t.ItemAmount)
		}
	}
```

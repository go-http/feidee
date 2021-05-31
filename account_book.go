package feidee

type AccountBook struct {
	Categories      []Category
	Stores          []IdName
	Members         []IdName
	Accounts        []IdName
	Projects        []IdName
	AccountInfoList map[int64]AccountInfo
}

//根据科目名获取科目ID为索引的Map
func (accountBook AccountBook) CategoryIdMap() map[int]Category {
	m := make(map[int]Category)
	for _, category := range accountBook.Categories {
		m[category.Id] = category
	}

	return m
}

//根据科目名获取科目ID
func (accountBook AccountBook) CategoryIdByName(name string) int {
	for _, item := range accountBook.Categories {
		if item.Name == name {
			return item.Id
		}
	}

	return 0
}

//根据商户名获取商户ID
func (accountBook AccountBook) StoreIdByName(name string) int {
	for _, item := range accountBook.Stores {
		if item.Name == name {
			return item.Id
		}
	}

	return 0
}

//根据成员名获取成员ID
func (accountBook AccountBook) MemberIdByName(name string) int {
	for _, item := range accountBook.Members {
		if item.Name == name {
			return item.Id
		}
	}

	return 0
}

//根据账户名获取账户ID
func (accountBook AccountBook) AccountIdByName(name string) int {
	for _, item := range accountBook.Accounts {
		if item.Name == name {
			return item.Id
		}
	}

	return 0
}

//根据项目名获取项目ID
func (accountBook AccountBook) ProjectIdByName(name string) int {
	for _, item := range accountBook.Projects {
		if item.Name == name {
			return item.Id
		}
	}

	return 0
}

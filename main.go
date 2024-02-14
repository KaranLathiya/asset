package main

import (
	"asset/dal"
	"asset/model"
	"asset/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

func main() {
	db, _ := dal.Connect()
	defer db.Close()
	fmt.Println("Server started")
	http.HandleFunc("/asset", asset)
	http.ListenAndServe(":8080", nil)
}

func asset(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("called")
	db := dal.GetDB()
	assetType := r.FormValue("asset_type")
	assetType = strings.TrimSpace(assetType)
	name := r.FormValue("name")
	name = strings.TrimSpace(name)
	var query string
	var filterArgsList []interface{}
	index := 1
	if name == "" && assetType == "" {
		query = "select asset_type, id, name, thumbnail_file_name, parent_id from public.asset_category order by parent_id asc, name asc"
	} else {
	var where []string
		if name != "" {
			where = append(where, "name ILIKE '%' || $"+strconv.Itoa(index)+" || '%'")
			filterArgsList = append(filterArgsList, name)
			index += 1
		}
		if assetType != "" {
			where = append(where, "asset_type = $"+strconv.Itoa(index))
			filterArgsList = append(filterArgsList, assetType)
		}
		query = fmt.Sprintf(
			`
		select
		asset_type,
		id,
		name,
		thumbnail_file_name,
		parent_id
	from
		public.asset_category
	where
		id in
		  (
		select
			id
		from
			public.asset_category
		where
			parent_id in
			    (
			select
				case
					when parent_id is null then id
				end as parent_id_null
			from
				public.asset_category
			where
				%v ))
		or id in(
		select
			case
				when parent_id is not null then parent_id
			end as parent_id_not_null
		from
			public.asset_category
		where
			%v)
		or id in(
		select
				id
		from
				public.asset_category
		where
				%v)
	order by
		parent_id asc,
		name asc;
		`, strings.Join(where, " AND "), strings.Join(where, " AND "), strings.Join(where, " AND "))
		query = sqlx.Rebind(sqlx.DOLLAR, query)
	}
	fmt.Println(query)
	rows, err := db.Query(query, filterArgsList...)
	if err != nil {
		databaseErrorMessage, databaseErrorCode := response.DatabaseErrorShow(err)
		response.MessageShow(databaseErrorCode, databaseErrorMessage, w)
		return
	}
	var mainCategoryMap = make(map[string]int)
	var assetTypeMap = make(map[string][]model.MainCategory)
	var mainCategoryList []model.MainCategory
	var parentID *string
	var i = 0
	for rows.Next() {
		var mainCategory model.MainCategory
		err = rows.Scan(&assetType, &mainCategory.ID, &mainCategory.Name, &mainCategory.Thumbnail, &parentID)
		if err != nil {
			response.MessageShow(400, err.Error(), w)
			return
		}
		if parentID == nil {
 			mainCategory.SubCategory = &[]model.MainCategory{}
			asset := assetTypeMap[assetType]
			asset = append(asset, mainCategory)
			assetTypeMap[assetType] = asset
			mainCategoryList = append(mainCategoryList, mainCategory)
			parentID := mainCategory.ID
			mainCategoryMap[parentID] = i
			i += 1
		} else {
			index := mainCategoryMap[*parentID]
			subCategory := model.MainCategory{
				ID:        mainCategory.ID,
				Name:      mainCategory.Name,
				Thumbnail: mainCategory.Thumbnail,
			  }
			*(mainCategoryList[index].SubCategory) = append(*(mainCategoryList[index].SubCategory), subCategory)

		}
	}
	fmt.Println(mainCategoryList)
	assetListData, _ := json.MarshalIndent(assetTypeMap, "", "  ")
	w.Write(assetListData)
}

package main

import (
	"asset/dal"
	"asset/model"
	"asset/response"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	db, _ := dal.Connect()
	defer db.Close()
	fmt.Println("Server started")
	http.HandleFunc("/asset", asset)
	http.ListenAndServe(":8080", nil)
}

func asset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called")
	db := dal.GetDB()
	name := r.FormValue("name")
	fmt.Println(name)
	rows, err := db.Query(
		`
		select
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
					name ilike '%' || $1 || '%'))
			or id in(
			select
				case
					when parent_id is not null then parent_id
				end as parent_id_not_null
			from
				public.asset_category
			where
				name ilike '%' || $1 || '%')
			or id in(
			select 
				id
			from
				public.asset_category
			where
				name ilike '%' || $1 || '%')
		order by
		parent_id  asc,
			"name" asc
				;
		`,name)
	if err != nil {
		databaseErrorMessage, databaseErrorCode := response.DatabaseErrorShow(err)
		response.MessageShow(databaseErrorCode, databaseErrorMessage, w)
		return
	}
	// fmt.Println(rows)
	var mainCategoryMap = make(map[string]int)
	var assetList []model.AssetList
	var parentID *string
	var i = 0
	for rows.Next() {
		var asset model.AssetList
		err = rows.Scan(&asset.ID, &asset.Name, &asset.Thumbnail, &parentID)
		if err != nil {
			response.MessageShow(400, err.Error(), w)
			return
		}
		if parentID == nil {
			assetList = append(assetList, asset)
			parentID := asset.ID 
			mainCategoryMap[parentID] = i
			i += 1
		} else {
			index := mainCategoryMap[*parentID]
			var subAsset model.SubAssetList
			subAsset.ID =asset.ID
			subAsset.Name =asset.Name
			subAsset.Thumbnail =asset.Thumbnail
			assetList[index].SubCategories = append(assetList[index].SubCategories, subAsset)
		}
	}
	var outputMap = make(map[string][]model.AssetList)
	outputMap["asset_type"] = assetList
	// fmt.Println(assetList)
	assetListData, _ := json.MarshalIndent(outputMap, "", "  ")
	w.Write(assetListData)
}

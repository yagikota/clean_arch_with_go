package outputdata

// response用のstruct
type CollectionsResponse struct {
	Collections []*CollectionResponse `json:"collections"`
}

// response用のstruct
type CollectionResponse struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

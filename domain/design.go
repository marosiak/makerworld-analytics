package domain

type DesignID int
type PublishedDesign struct {
	ID   DesignID `json:"designId"`
	Name string   `json:"name"`
}

type PublishedDesignsList []PublishedDesign

func (p PublishedDesignsList) Exists(id DesignID) bool {
	for _, design := range p {
		if design.ID == id {
			return true
		}
	}
	return false
}

package data

import "time"

type Property struct {
	ID        int64
	CreatedAt time.Time
	Version   int32
	County    string
	PropID    string
	AltID     string
	Owners    []string

	// Addressing
	SitusAddress string
	OwnerAddress string
	City         string
	State        string
	ZIP          string

	// Classification
	Class        string
	Neighborhood string
	Acres        float64
	Type         string
	Grade        string
	YearBuilt    int16

	// Building Attributes
	Occupancy     string
	RoofStruct    string
	RoofCover     string
	Heating       string
	AC            string
	Stories       float32
	Bedrooms      int16
	Bathrooms     float32
	ExteriorWall  string
	InteriorFloor string

	// Value History (latest or summary)
	ValueYear        int32
	ValueReason      string
	LandValue        int64
	ImprovementValue int64
	TotalAppraised   int64
	TotalAssessed    int64

	// Transfer History (latest or summary)
	TransferBook string
	TransferPage string
	TransferDate time.Time
	Grantor      string
	Grantee      string
	DeedType     string
	VacantLand   bool
	SalePrice    int64

	// Land & Features
	LandPrimaryUse string
	LandType       string
	EffFrontage    int32
	EffDepth       int32
	FloorAreaCode  string
	FloorDesc      string
	GrossAreaSF    int32
	FinishedSF     int32
	FloorConstr    string
	ExtFeatCode    string
	ExtFeatDesc    string
	ExtFeatSizeSF  int32
	ExtFeatConstr  string

	// Legal / References
	LegalDescription []string
	SourceURL        string
	Photos           []string
}

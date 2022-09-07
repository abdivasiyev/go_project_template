package models

type GetDashboardResponse struct {
	DriverStatistics  GetDriverDashboardResponse  `json:"driver_statistics"`
	TruckStatistics   GetTruckDashboardResponse   `json:"truck_statistics"`
	TrailerStatistics GetTrailerDashboardResponse `json:"trailer_statistics"`
	CarStatistics     GetCarDashboardResponse     `json:"car_statistics"`
}

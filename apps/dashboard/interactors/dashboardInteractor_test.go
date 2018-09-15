package interactors

import (
	"testing"
	"stock/apps/dashboard/interfaces"
)

func BenchmarkDashboardInteractor_GetActivityLog(b *testing.B) {
	var d interfaces.DashboardUseCases
	d = DashboardInteractor{}
	for i:=0; i < b.N ; i++{
		d.GetActivityLog("0,1536479566",0)
	}

}
package providers

import (
	"strconv"

	"github.com/Dmk88/go_practice/currencymonitor/models"
	"github.com/gocql/gocql"
)

type monitoringProvider struct {
	scylla *gocql.Session
}

func NewMonitoringProvider(scylla *gocql.Session) monitoringProvider {
	return monitoringProvider{scylla: scylla}
}

// TODO: create query builder
func (p *monitoringProvider) FindByID(id string) (*models.Monitoring, error) {
	monitoring := new(models.Monitoring)
	err := p.scylla.Query(`SELECT id, status, period, frequency, data
									FROM currencymonitor.monitoring
									WHERE id=? LIMIT 1`, id).
		Scan(&monitoring.ID, &monitoring.Status, &monitoring.Period, &monitoring.Frequency, &monitoring.Data)

	return monitoring, err
}

func (p *monitoringProvider) FindByStatus(monitoringStatus int) []models.Monitoring {
	monitoringRecords := make([]models.Monitoring, 0)

	iter := p.scylla.Query(`SELECT id, status, period, frequency, data
									FROM currencymonitor.monitoring
									WHERE status=?`, monitoringStatus).Iter()

	var monitoring models.Monitoring
	for iter.Scan(&monitoring.ID, &monitoring.Status, &monitoring.Period, &monitoring.Frequency, &monitoring.Data) {
		monitoringRecords = append(monitoringRecords, monitoring)
	}

	return monitoringRecords
}

func (p *monitoringProvider) Save(monitoring *models.Monitoring, ttl int) error {
	return p.scylla.Query(`INSERT INTO currencymonitor.monitoring(id, status, period, frequency, data)
    						VALUES (?,?,?,?,?) USING TTL `+strconv.Itoa(ttl),
		monitoring.ID, monitoring.Status, monitoring.Period, monitoring.Frequency, monitoring.Data).Exec()
}

func (p *monitoringProvider) UpdateData(monitoringId, monitoringData string) error {
	return p.scylla.Query(`UPDATE currencymonitor.monitoring
							SET data = ? WHERE id = ?`, monitoringData, monitoringId).Exec()
}

func (p *monitoringProvider) UpdateStatus(monitoringId string, monitoringStatus int) error {
	return p.scylla.Query(`UPDATE currencymonitor.monitoring
							SET status = ? WHERE id = ?`, monitoringStatus, monitoringId).Exec()
}

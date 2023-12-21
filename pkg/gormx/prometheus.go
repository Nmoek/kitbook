package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	vector *prometheus.SummaryVec
}

func (c *Callbacks) Name() string {
	return "prometheus"
}

func (c *Callbacks) Initialize(db *gorm.DB) error {
	/*create*/
	err := db.Callback().Create().Before("*").Register("prometheus_gorm_create_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Create().After("*").Register("prometheus_gorm_create_after", c.After("CREATE"))
	if err != nil {
		return err
	}

	/*query*/
	err = db.Callback().Query().Before("*").Register("prometheus_gorm_query_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("*").Register("prometheus_gorm_query_after", c.After("QUERY"))
	if err != nil {
		return err
	}

	/*raw*/
	err = db.Callback().Raw().Before("*").Register("prometheus_gorm_raw_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Raw().After("*").Register("prometheus_gorm_query_after", c.After("RAW"))
	if err != nil {
		return err
	}

	/*update*/
	err = db.Callback().Update().Before("*").Register("prometheus_gorm_update_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("*").Register("prometheus_gorm_update_after", c.After("UPDATE"))
	if err != nil {
		return err
	}

	/*delete*/
	err = db.Callback().Delete().Before("*").Register("prometheus_gorm_delete_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("*").Register("prometheus_gorm_delete_after", c.After("DELETE"))
	if err != nil {
		return err
	}

	/*row*/
	err = db.Callback().Row().Before("*").Register("prometheus_gorm_row_before", c.Before())
	if err != nil {
		return err
	}

	err = db.Callback().Row().After("*").Register("prometheus_gorm_row_after", c.After("ROW"))

	return err

}

func NewCallbacks(opts prometheus.SummaryOpts) *Callbacks {
	vector := prometheus.NewSummaryVec(opts, []string{"type", "table"})
	prometheus.MustRegister(vector)
	return &Callbacks{
		vector: vector,
	}
}

func (c *Callbacks) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}

func (c *Callbacks) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		start, ok := val.(time.Time)
		if ok {
			// us微秒
			duration := time.Since(start).Microseconds()
			c.vector.WithLabelValues(typ, db.Statement.Table).Observe(float64(duration))
		}
	}
}

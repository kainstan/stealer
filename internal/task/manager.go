package task

import (
	"go.uber.org/zap"
	"sync"
	"stealer/log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type ZapLogger struct {
	cron.Logger

	log *zap.SugaredLogger
}

func (z *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.log.Info(msg, keysAndValues)
}

func (z *ZapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	z.log.Error(err, msg, keysAndValues)
}

// https://www.cnblogs.com/jssyjam/p/11910851.html
type CronManager struct {

	inner *cron.Cron
	ids   map[string]cron.EntryID
	mutex sync.Mutex
}

var (
	crontab *CronManager
	lock sync.Mutex
)

// New init
func New() (s *CronManager) {
	cronExpr := cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	))

	s = &CronManager{
		inner: cron.New(
			//cron.WithSeconds(),
			cronExpr,
			cron.WithChain(cron.Recover(&ZapLogger{log: log.AppLog})),
		),
		//inner: cron.New(cron.WithParser(cron.NewParser(
		//	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		//))),
		ids:   make(map[string]cron.EntryID),
	}
	return s
}

//func init()  {
//	crontab = New()
//	crontab.Start()
//}

func CronManger() *CronManager {
	if crontab != nil {
		return crontab
	}
	lock.Lock()
	if crontab != nil {
		return crontab
	}
	crontab = New()
	crontab.Start()
	defer lock.Unlock()
	return crontab
}

// IDs ...
func (c *CronManager) IDs() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	validIDs := make([]string, 0, len(c.ids))
	invalidIDs := make([]string, 0)
	for sid, eid := range c.ids {
		if e := c.inner.Entry(eid); e.ID != eid {
			invalidIDs = append(invalidIDs, sid)
			continue
		}
		validIDs = append(validIDs, sid)
	}
	for _, id := range invalidIDs {
		delete(c.ids, id)
	}
	return validIDs
}

// Start start the crontab engine
func (c *CronManager) Start() {
	c.inner.Start()
}

// Stop stop the crontab engine
func (c *CronManager) Stop() {
	c.inner.Stop()
}

// DelByID remove one crontab task
func (c *CronManager) DelByID(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	eid, ok := c.ids[id]
	if !ok {
		return
	}
	c.inner.Remove(eid)
	delete(c.ids, id)
}

// AddByID add one crontab task
// id is unique
// spec is the crontab expression
func (c *CronManager) AddByID(id string, spec string, cmd cron.Job) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.ids[id]; ok {
		return errors.Errorf("crontab id exists")
	}
	eid, err := c.inner.AddJob(spec, cmd)
	if err != nil {
		return err
	}
	c.ids[id] = eid
	return nil
}

// AddByFunc add function as crontab task
func (c *CronManager) AddByFunc(id string, spec string, f func()) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.ids[id]; ok {
		return errors.Errorf("crontab id exists")
	}
	eid, err := c.inner.AddFunc(spec, f)
	if err != nil {
		return err
	}
	c.ids[id] = eid
	return nil
}

// IsExists check the crontab task whether existed with job id
func (c *CronManager) IsExists(jid string) bool {
	_, exist := c.ids[jid]
	return exist
}

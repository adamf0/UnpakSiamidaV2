package helper

import "time"

type DateContext struct {
	strategy IDateFormat
}

func (c *DateContext) SetStrategy(s IDateFormat) {
	c.strategy = s
}

func (c *DateContext) NameDay(t time.Time) string {
	return c.strategy.NameDay(t)
}

func (c *DateContext) Day(t time.Time) string {
	return c.strategy.Day(t)
}

func (c *DateContext) Month(t time.Time) string {
	return c.strategy.Month(t)
}

func (c *DateContext) Year(t time.Time) string {
	return c.strategy.Year(t)
}

func (c *DateContext) Format(t time.Time) string {
	return c.strategy.Format(t)
}

func (c *DateContext) FormatWithDay(t time.Time) string {
	return c.strategy.FormatWithDay(t)
}

func (c *DateContext) FormatWithTime(t time.Time) string {
	return c.strategy.FormatWithTime(t)
}

func (c *DateContext) FormatDefault(t *time.Time) string {
	return c.strategy.FormatDefault(t)
}

func (c *DateContext) FormatWithDayDefault(t *time.Time) string {
	return c.strategy.FormatWithDayDefault(t)
}

func (c *DateContext) FormatWithTimeDefault(t *time.Time) string {
	return c.strategy.FormatWithTimeDefault(t)
}

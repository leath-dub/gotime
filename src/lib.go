package lib

import "time"

type BodyTemplate struct {
    FirstDayInWeek string
    Name string
    DayOfWeek int
    CategoryIdentities []string
}

func startOfWeek(date time.Time) time.Time {
    offset := int(date.Weekday()) - 1
    return date.AddDate(0, 0, -offset)
}

func NewBody(date time.Time, id []string) BodyTemplate {
    start := startOfWeek(date)

    return BodyTemplate {
        FirstDayInWeek: start.Format(time.RFC3339),
        Name: date.Weekday().String(),
        DayOfWeek: int(date.Weekday()),
        CategoryIdentities: id,
    }
}

type CategoryEvent struct {
    EventIdentity string
    HostKey string
    Description string
    EndDateTime string
    EventType string
    IsPublished string
    Location string
    Owner string
    StartDateTime string
    IsDeleted bool
    LastModified string
    ExtraProperties []struct {
        Name string
        DisplayName string
        Value string
        Rank int
    }
}

type CategoryEvents []CategoryEvent

type ResponseTemplate struct {
    CategoryTypeIdentity string
    CategoryTypeName string
    CategoryEvents CategoryEvents
}

func getMinutes(time *time.Time) int {
    return (time.Hour() * 60) + time.Minute()
}

/* Functions Len, Less & Swap setup a sorting interface (sort.Interface) */
func (self CategoryEvents) Len() int {
    return len(self)
}

func (self CategoryEvents) Less(i, j int) bool {
    /* Parse date into Time type */
    iTime, _ := time.Parse(time.RFC3339, self[i].StartDateTime)
    jTime, _ := time.Parse(time.RFC3339, self[j].StartDateTime)

    /* get time in minutes */
    return getMinutes(&iTime) < getMinutes(&jTime)
}


func (self CategoryEvents) Swap(i, j int) {
    tmp := self[i]
    self[i] = self[j]
    self[j] = tmp
}

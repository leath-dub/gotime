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
/*
type BodyTemplate struct {
    ViewOptions struct {
        Days []struct {
            Name string
            DayOfWeek int
            IsDefault bool
        }
        Weeks []struct {
            WeekNumber int
            WeekLabel int
            FirstDayInWeek string
        }
        TimePeriods []struct {
            Description string
            StartTime string
            EndTime string
            IsDefault bool
        }
        DatePeriods []struct {
            Description string
            StartDateTime string
            EndDateTime string
            IsDefault bool
            IsThisWeek bool
            IsNextWeek bool
            Type string
        }
        LegendItems []any
        InstitutionConfig struct {}
        DateConfig struct {
          FirstDayInWeek int
          StartDate string
          EndDate string
        }
    }
    CategoryIdentities []string
}
*/

type ResponseTemplate struct {
    CategoryTypeIdentity string
    CategoryTypeName string
    CategoryEvents []struct {
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
}

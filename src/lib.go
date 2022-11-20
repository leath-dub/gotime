package lib

type BodyTemplate struct {
    FirstDayInWeek string
    Name string
    DayOfWeek int
    CategoryIdentities []string
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


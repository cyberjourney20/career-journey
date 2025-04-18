package models

import (
	"time"
)

// User is the user model
type User struct {
	UserID      string
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Skills      []UserSkills
	Certs       []UserCerts
	Preferences []UserPreferences
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Contact struct {
	ID              int
	UserID          string
	FirstName       string
	LastName        string
	Email           string
	CompanyID       int
	JobTitle        string
	MobilePhone     string
	WorkPhone       string
	Phone3          string
	Linkedin        string
	Github          string
	Website         string
	Tags            []int
	Notes           string
	Description     string
	Objective       string
	ContactTimeLine string
	Favorite        bool
	Company         Company
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Application struct {
	ID          int
	UserID      string
	Company     Company
	JobID       string
	POC         Contact
	AppStatus   string
	Notes       string
	AppTimeLine string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type JobListing struct {
	ID             int
	ExternalID     string            `json:"job_id"`
	CompanyID      int               `json:"company_id"`
	URL            string            `json:"url"`
	JobTitle       string            `json:"job_title"`
	JobDescription string            `json:"job_description"`
	WorkSetting    int               `json:"work_setting"`
	ReqYOE         int               `json:"req_yoe"`
	ReqSkills      []ReferenceSkills `json:"req_skills"`
	ReqCerts       []ReferenceCerts  `json:"req_certs"`
	LowPay         string            `json:"low_pay"`
	HighPay        string            `json:"high_pay"`
	TargetPay      string            `json:"target_pay"`
	Location       Location          `json:"location"`
	Company        Company           `json:"company"`
	Posted         time.Time
	Closes         time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Company struct {
	ID          int
	CompanyName string `json:"company_name"`
	URL         string
	Address     Location
	Industry    string
	Size        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Location struct {
	ID      int
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Zip     string
}

type Tag struct {
	TagID int
}

type POC struct {
	ContactID int
}

type TimeLineData struct {
	Marker time.Time
	Event  string
}

type UserSkills struct {
	ReferenceSkillsID int
	Current           bool
	StartDate         time.Time
	EndDate           time.Time
	SkillLevel        int
}

type UserCerts struct {
	ReferenceCerts int
	Expires        bool
	AwardDate      time.Time
	ExpiryDate     time.Time
}

type UserPreferences struct {
	ShowTour  bool
	DarkMode  bool
	TargetPay int
}

type ReferenceSkills struct {
	ID            int
	Domain        string
	Title         string
	Description   string
	RelatedSkills []RelatedSkills
	Trend         []Trend
}

type ReferenceCerts struct {
	ID            int
	CertBody      string
	Url           string
	Title         string
	Cost          int
	Description   string
	RelatedSkills []RelatedSkills
	Trend         []Trend
}

type RelatedSkills struct {
	ReferenceSkillsID int
	Strength          int
}

type Trend struct {
	StartDate time.Time
	EndDate   time.Time
	Count     int
}

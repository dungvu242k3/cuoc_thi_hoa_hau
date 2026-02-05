package domain

import "time"

type ContestantStatus string

const (
	ContestantStatusDraft    ContestantStatus = "draft"
	ContestantStatusPending  ContestantStatus = "pending"
	ContestantStatusApproved ContestantStatus = "approved"
	ContestantStatusRejected ContestantStatus = "rejected"
	ContestantStatusDeleted  ContestantStatus = "deleted"
)

type Contestant struct {
	ID       string           `bson:"_id,omitempty" json:"id"`
	UserID   string           `bson:"user_id" json:"userId"`
	SBD      string           `bson:"sbd" json:"sbd"` // Số báo danh
	Status   ContestantStatus `bson:"status" json:"status"`
	IsPublic bool             `bson:"is_public" json:"isPublic"`

	PersonalInfo PersonalInfo   `bson:"personal_info" json:"personalInfo"`
	PhysicalInfo PhysicalInfo   `bson:"physical_info" json:"physicalInfo"`
	SkillEdu     SkillEducation `bson:"skill_edu" json:"skillEducation"`
	Portfolio    Portfolio      `bson:"portfolio" json:"portfolio"`
	VoteCount    int64          `bson:"vote_count" json:"voteCount"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

type PersonalInfo struct {
	FullName     string    `bson:"full_name" json:"fullName"`
	DateOfBirth  time.Time `bson:"dob" json:"dob"`
	Nationality  string    `bson:"nationality" json:"nationality"`
	Gender       string    `bson:"gender" json:"gender"`
	IdentityCard string    `bson:"identity_card" json:"identityCard"`
	Phone        string    `bson:"phone" json:"phone"`
	Email        string    `bson:"email" json:"email"`
	Address      string    `bson:"address" json:"address"`
	Job          string    `bson:"job" json:"job"`
}

type PhysicalInfo struct {
	Height       float64 `bson:"height" json:"height"`
	Weight       float64 `bson:"weight" json:"weight"`
	Measurements string  `bson:"measurements" json:"measurements"` // "90-60-90"
}

type SkillEducation struct {
	EducationLevel string   `bson:"education_level" json:"educationLevel"`
	Languages      []string `bson:"languages" json:"languages"`
	Skills         []string `bson:"skills" json:"skills"`
}

type Portfolio struct {
	AvatarURL    string   `bson:"avatar_url" json:"avatarUrl"`
	GalleryURLs  []string `bson:"gallery_urls" json:"galleryUrls"`
	VideoURL     string   `bson:"video_url" json:"videoUrl"`
	Introduction string   `bson:"introduction" json:"introduction"`
	SocialLinks  []string `bson:"social_links" json:"socialLinks"`
}

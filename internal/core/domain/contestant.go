package domain

import (
	"time"
)

// trạng thái hồ sơ
type ContestantStatus string

const (
	StatusDraft    ContestantStatus = "DRAFT"    // bản nháp
	StatusPending  ContestantStatus = "PENDING"  // chờ duyệt
	StatusApproved ContestantStatus = "APPROVED" // đa phê duyệt
	StatusRejected ContestantStatus = "REJECTED" // bị từ chối
	StatusLocked   ContestantStatus = "LOCKED"   // tài khoản bị khóa ( gian lận ...)
	StatusDeleted  ContestantStatus = "DELETED"  // đã xóa mềm (soft delete)
)

// Struct chính
type Contestant struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	UserID       string `bson:"user_id" json:"userId"`
	IdentifyCard string `bson:"identify_card" json:"identifyCard"`
	SBD          string `bson:"sbd" json:"sbd"` // Số báo danh

	PersonalInfo PersonalInfo   `bson:"personal_info" json:"personalInfo"`
	PhysicalInfo PhysicalInfo   `bson:"physical_info" json:"physicalInfo"`
	SkillEdu     SkillEducation `bson:"skill_education" json:"skillEducation"`
	Portfolio    Portfolio      `bson:"portfolio" json:"portfolio"`

	Status   ContestantStatus `bson:"status" json:"status"`
	Reason   string           `bson:"reason,omitempty" json:"reason"`
	IsPublic bool             `bson:"is_public" json:"isPublic"`

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// Các struct thành phần
type PersonalInfo struct {
	FullName    string    `bson:"full_name" json:"fullName"`
	DateOfBirth time.Time `bson:"dob" json:"dob"`
	Nationality string    `bson:"nationality" json:"nationality"`
	Job         string    `bson:"job" json:"job"` // Nghề nghiệp
	Phone       string    `bson:"phone" json:"phone"`
	Email       string    `bson:"email" json:"email"`
	Address     string    `bson:"address" json:"address"`
}

type PhysicalInfo struct {
	Height       float64 `bson:"height" json:"height"`
	Weight       float64 `bson:"weight" json:"weight"`
	Measurements string  `bson:"measurements" json:"measurements"`
}

type SkillEducation struct {
	EducationLevel string   `bson:"education_level" json:"educationLevel"`
	Languages      []string `bson:"languages" json:"languages"`
	Skills         []string `bson:"skills" json:"skills"`
}

type Portfolio struct {
	AvatarURL    string   `bson:"avatar_url" json:"avatarUrl"`
	GalleryURLs  []string `bson:"gallery_urls" json:"galleryUrls"`
	Introduction string   `bson:"introduction" json:"introduction"`
	SocialLinks  []string `bson:"social_links" json:"socialLinks"`
}

// Helper để tạo bản sao an toàn (Xóa thông tin nhạy cảm)
func (c *Contestant) ToPublicView() *Contestant {
	safe := *c
	// Xóa thông tin định danh hệ thống & cá nhân nhạy cảm
	safe.UserID = ""
	safe.IdentifyCard = ""
	safe.PersonalInfo.Phone = ""
	safe.PersonalInfo.Email = ""
	safe.PersonalInfo.Address = ""

	if !c.IsPublic {
		safe.Portfolio.GalleryURLs = nil
	}

	// Public View: Giữ lại SBD và Job
	// safe.SBD = c.SBD (Mặc định đã copy)
	// safe.PersonalInfo.Job = c.PersonalInfo.Job (Mặc định đã copy)

	return &safe
}

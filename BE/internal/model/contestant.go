/*
Package model định nghĩa Cấu trúc Dữ liệu (Entities/Domain Models) cốt lõi của ứng dụng.
Tác dụng của contestant.go:
- Ánh xạ 1-1 với cấu trúc tài liệu (document) lưu trong bảng "contestants" của MongoDB.
- Chứa các thẻ (tags) như `bson:` để MongoDB biết cách đọc/ghi, và `json:` để Frontend nhận JSON khớp tên.
- Không chứa logic nghiệp vụ, thuần túy là định dạng dữ liệu lưu trữ.
*/
package model

import "time"

type ContestantStatus string

const (
	ContestantStatusDraft    ContestantStatus = "draft"
	ContestantStatusPending  ContestantStatus = "pending"
	ContestantStatusApproved ContestantStatus = "approved"
	ContestantStatusRejected ContestantStatus = "rejected"
	ContestantStatusDeleted  ContestantStatus = "deleted"
)

// Contestant là cấu trúc lưu trữ toàn bộ hồ sơ đăng ký dự thi của một vòng thi.
// Nó chia nhỏ thông tin cá nhân, hình thể, học vấn thành các struct lồng nhau (embedded structs) để dễ quản lý.
type Contestant struct {
	ID       string           `bson:"_id,omitempty" json:"id"` // ObjectID do Mongo tự sinh
	UserID   string           `bson:"user_id" json:"userId"`   // Khóa ngoại liên kết với bảng Users (Tài khoản)
	SBD      string           `bson:"sbd" json:"sbd"`          // Số báo danh
	Status   ContestantStatus `bson:"status" json:"status"`
	IsPublic bool             `bson:"is_public" json:"isPublic"` // Cờ hiệu đánh dấu hồ sơ có được hiển thị lên trang chủ không

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

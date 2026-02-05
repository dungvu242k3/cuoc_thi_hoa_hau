export const DASHBOARD_QUERY = `
  query GetDashboardData {
    publicContestants(limit: 100) {
      id
      sbd
      personalInfo {
        fullName
      }
      portfolio {
        avatarUrl
      }
      status
      isPublic
    }
    myScores {
      sbd
      totalScore
    }
  }
`;

export const ADMIN_QUERY = `
  query GetPendingContestants {
    adminContestants(status: "pending", limit: 50) {
      id
      sbd
      personalInfo {
        fullName
        dob
        nationality
        gender
        identityCard
        phone
        email
        address
        job
      }
      physicalInfo {
        height
        weight
        measurements
      }
      skillEducation {
        educationLevel
        languages
        skills
      }
      portfolio {
        avatarUrl
        galleryUrls
        videoUrl
        introduction
        socialLinks
      }
      createdAt
    }
  }
`;

export const APPROVE_MUTATION = `
  mutation ApproveContestant($id: ID!, $isApproved: Boolean!) {
    approveContestant(id: $id, isApproved: $isApproved) {
      id
      status
      isPublic
    }
  }
`;

export const GET_CONTESTANT = `
  query GetContestant($id: ID!) {
     publicContestantDetail(id: $id) {
        id sbd
        personalInfo {
            fullName dob hometown: address
            gender nationality job phone email
        }
        physicalInfo {
            height weight measurements
        }
        skillEducation {
            educationLevel languages skills
        }
        portfolio {
            avatarUrl galleryUrls videoUrl introduction socialLinks
        }
     }
  }
`;

export const SUBMIT_SCORE = `
  mutation SubmitScore($input: ScoreInput!) {
    submitScore(input: $input) {
      id
      totalScore
    }
  }
`;

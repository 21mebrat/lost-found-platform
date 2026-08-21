package validator

import domain "github.com/21mebrat/lost-found-platform/internal/domain/otp"

func IsValidPurpose(
	purpose domain.Purpose,
) bool {
	switch purpose {
	case domain.PurposeLogin,
		domain.PurposeRegistration,
		domain.PurposePasswordReset,
		domain.PurposePhoneChange:
		return true

	default:
		return false
	}
}

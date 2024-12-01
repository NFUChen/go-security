package service

import (
	"context"
	. "go-security/erp/internal/repository"
)

type ProfilePricingService struct {
	ProfileService       *ProfileService
	PricingPolicyService *PricingPolicyService
}

func (service *ProfilePricingService) ApplyPricingPolicyToProfile(ctx context.Context, profileID uint, policyID uint) error {
	profile, err := service.ProfileService.GetProfileByID(ctx, profileID)
	if err != nil {
		return err
	}
	policy, err := service.PricingPolicyService.GetPolicyByID(ctx, policyID)
	if err != nil {
		return err
	}
	profile.PricingPolicy = *policy
	_, err = service.ProfileService.UpdateProfile(ctx, profile.UserID, &UserProfile{PricingPolicyID: policyID})
	if err != nil {
		return err // Handle update failure
	}

	return nil
}

func NewProfilePricingService(profileService *ProfileService, pricingPolicyService *PricingPolicyService) *ProfilePricingService {
	return &ProfilePricingService{
		ProfileService:       profileService,
		PricingPolicyService: pricingPolicyService,
	}
}

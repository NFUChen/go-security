package service

import "context"

type ProfilePricingService struct {
	ProfileService       *ProfileService
	PricingPolicyService *PricingPolicyService
}

func (service *ProfilePricingService) ApplyPricingPolicyToProfile(ctx context.Context, profileID uint, policyID uint) error {
	profile, err := service.ProfileService.GetProfileByID(ctx, profileID)
	if err != nil {
		return err
	}
	policy, err := service.PricingPolicyService.GetPolicyByID(policyID)
	if err != nil {
		return err
	}
	profile.PricingPolicy = *policy
	err = service.ProfileService.UpdateProfile(ctx, profile, map[string]any{"pricing_policy_id": policyID})
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

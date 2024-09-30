package ops

import (
	"context"
	"database/sql"

	"github.com/friendlycaptcha/friendly-stripe-sync/db/postgres"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/utils"
	"github.com/stripe/stripe-go/v74"
)

func (o *Ops) HandleCouponUpdated(c context.Context, partialCoupon *stripe.Coupon) error {
	coup := partialCoupon
	if partialCoupon.AppliesTo == nil {
		params := &stripe.CouponParams{}
		params.AddExpand("applies_to")
		var err error
		coup, err = o.stripe.Coupons.Get(coup.ID, params)
		if err != nil {
			return err
		}
	}

	return o.db.Q.UpsertCoupons(c, postgres.UpsertCouponsParams{
		ID:               coup.ID,
		Object:           coup.Object,
		AmountOff:        sql.NullInt64{Int64: coup.AmountOff, Valid: coup.AmountOff != 0},
		Created:          coup.Created,
		Currency:         utils.StringToNullString(string(coup.Currency)),
		Duration:         string(coup.Duration),
		DurationInMonths: sql.NullInt64{Int64: coup.DurationInMonths, Valid: coup.DurationInMonths != 0},
		MaxRedemptions:   sql.NullInt64{Int64: coup.MaxRedemptions, Valid: coup.MaxRedemptions != 0},
		Metadata:         utils.MarshalToNullRawMessage(coup.Metadata),
		Name:             coup.Name,
		PercentOff:       sql.NullFloat64{Float64: coup.PercentOff, Valid: coup.PercentOff != 0},
		RedeemBy:         sql.NullInt64{Int64: coup.RedeemBy, Valid: coup.RedeemBy != 0},
		TimesRedeemed:    sql.NullInt64{Int64: coup.TimesRedeemed, Valid: coup.TimesRedeemed != 0},
		AppliesTo:        utils.MarshalToNullRawMessage(coup.AppliesTo),
		Valid:            coup.Valid,
	})
}

func (o *Ops) HandleCouponDeleted(c context.Context, coupon *stripe.Coupon) error {
	return o.db.Q.DeleteCoupon(c, coupon.ID)
}

func (o *Ops) EnsureCouponLoaded(c context.Context, couponID string) error {
	exists, err := o.db.Q.CouponExists(c, couponID)
	if err != nil {
		return err
	}
	if !exists {
		// HandleCouponUpdated will fetch the coupon from Stripe because AppliesTo isn't set
		o.HandleCouponUpdated(c, &stripe.Coupon{ID: couponID})
	}

	return nil
}

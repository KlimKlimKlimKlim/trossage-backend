package controller

import (
	"context"
	"fmt"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres"
)

func (c *Controller) createTokens(ctx context.Context, tx postgres.IRepository, userID int64) (string, string, error) {
	accessString, err := c.accessJWTController.GenerateSignedToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshString, storeToken, err := c.refreshJWTController.GenerateSignedTokenAndModel(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	if _, err = tx.InsertToken(ctx, storeToken); err != nil {
		return "", "", fmt.Errorf("failed to insert token: %w", err)
	}

	return accessString, refreshString, nil
}

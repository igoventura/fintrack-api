package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/igoventura/fintrack-api/domain"
)

type TransactionService struct {
	repo        domain.TransactionRepository
	accountRepo domain.AccountRepository
}

func NewTransactionService(repo domain.TransactionRepository, accountRepo domain.AccountRepository) *TransactionService {
	return &TransactionService{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *TransactionService) List(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}
	return s.repo.List(ctx, tenantID, filter)
}

func (s *TransactionService) Create(ctx context.Context, t *domain.Transaction, tagIDs []string) error {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}
	t.TenantID = tenantID

	userID := domain.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}
	t.CreatedBy = userID
	t.UpdatedBy = userID

	// Validate basic fields
	if valid, errs := t.IsValid(); !valid {
		var errMsg string
		for field, err := range errs {
			errMsg += fmt.Sprintf("%s: %s; ", field, err.Error())
		}
		return errors.New("validation failed: " + errMsg)
	}

	// 1. Currency Logic: Default to FromAccount currency if not set
	if t.Currency == "" {
		fromAccount, err := s.accountRepo.GetByID(ctx, t.FromAccountID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to fetch from_account for currency: %w", err)
		}

		if fromAccount.TenantID != tenantID {
			return errors.New("from_account does not belong to this tenant")
		}
		t.Currency = fromAccount.Currency
	}

	// 2. ToAccount Validation (if applicable)
	if t.ToAccountID != nil && *t.ToAccountID != "" {
		toAccount, err := s.accountRepo.GetByID(ctx, *t.ToAccountID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to fetch to_account: %w", err)
		}
		if toAccount.TenantID != tenantID {
			return errors.New("to_account does not belong to this tenant")
		}
	}

	// 3. Create Transaction
	if err := s.repo.Create(ctx, t); err != nil {
		return err
	}

	// 4. Tags Association
	if len(tagIDs) > 0 {
		if err := s.repo.AddTagsToTransaction(ctx, t.ID, tagIDs); err != nil {
			return fmt.Errorf("transaction created but failed to link tags: %w", err)
		}
	}

	return nil
}

func (s *TransactionService) Update(ctx context.Context, t *domain.Transaction, tagIDs []string) error {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}
	t.TenantID = tenantID // Ensure we don't overwrite with wrong tenant

	userID := domain.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}
	t.UpdatedBy = userID

	if err := s.repo.Update(ctx, t); err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) error {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}
	userID := domain.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}
	return s.repo.Delete(ctx, tenantID, id, userID)
}

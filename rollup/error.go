// Copyright (C) 2023 Wooyang2018
// Licensed under the GNU General Public License v3.0

package rollup

import "errors"

var (
	// ErrSizeByteSlice memory checking
	ErrSizeByteSlice = errors.New("byte slice size is inconsistant with Account size")

	// ErrNonExistingAccount account not in the database
	ErrNonExistingAccount = errors.New("the account is not in the rollup database")

	// ErrWrongSignature wrong signature
	ErrWrongSignature = errors.New("invalid signature")

	// ErrAmountTooHigh the amount is bigger than the balance
	ErrAmountTooHigh = errors.New("amount is bigger than balance")

	// ErrNonce inconsistant nonce between transfer and account
	ErrNonce = errors.New("incorrect nonce")

	// ErrIndexConsistency the map publicKey(string) -> index(int) gives acces to the account position.
	// Account has a field index, that should match position.
	ErrIndexConsistency = errors.New("account's position should match account's index")
)

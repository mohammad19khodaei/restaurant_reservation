package utils_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := faker.Password()

	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isValid := utils.IsHashPasswordValid(hashedPassword, password)
	require.NoError(t, err)
	require.True(t, isValid)

	wrongPassword := faker.Password()
	isValid = utils.IsHashPasswordValid(wrongPassword, password)
	require.False(t, isValid)
}

package hasher

import (
	"runtime"

	"github.com/alexedwards/argon2id"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
)

const (
	maxUint8 = 255
)

type Hasher struct {
	params *argon2id.Params
}

func New(cfg *config.Hasher) *Hasher {
	return &Hasher{
		params: &argon2id.Params{
			Memory:      cfg.Memory,
			Iterations:  cfg.Iterations,
			Parallelism: getArgon2Parallelism(),
			SaltLength:  cfg.SaltLength,
			KeyLength:   cfg.KeyLength,
		},
	}
}

func getArgon2Parallelism() uint8 {
	cpus := runtime.NumCPU()
	if cpus > maxUint8 {
		return maxUint8
	}

	if cpus < 1 {
		return 1
	}

	return uint8(cpus)
}

func (h *Hasher) HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, h.params)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (h *Hasher) VerifyPassword(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

package cryptor

import "golang.org/x/crypto/bcrypt"

type Cryptor interface {
	EncryptKeyword(keyword string) (string, error)
	CompareHashAndKeyword(hash, keyword string) error
}

type bcryptor struct {
	pool *workerPool
}

func New(asyncHashingLimit int) Cryptor {
	return &bcryptor{
		pool: NewWorkerPool(asyncHashingLimit),
	}
}

func (c *bcryptor) EncryptKeyword(keyword string) (string, error) {
	resChan := make(chan string, 1)
	errChan := make(chan error, 1)

	c.pool.Add(func() {
		hash, err := bcrypt.GenerateFromPassword([]byte(keyword), bcrypt.MinCost)
		if err != nil {
			errChan <- err
			return
		}
		resChan <- string(hash)
	})

	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return "", err
	}
}

func (c *bcryptor) CompareHashAndKeyword(hash, keyword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(keyword))
}

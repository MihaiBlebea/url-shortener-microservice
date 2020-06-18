package shortener

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
)

// Errors
var (
	ErrRedirectNotFound = errors.New("Could not find redirect")
	ErrRedirectInvalid  = errors.New("Redirect is invalid")
)

// RedirectService is the main service forshortener
type redirectService struct {
	redirectRepository RedirectRepository
}

// NewRedirectService instantiates redirectService
func NewRedirectService(redirectRepository RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepository,
	}
}

// Find a redirect by code and increments it's hits counter
func (rs *redirectService) Find(code string) (redirect Redirect, err error) {
	redirect, err = rs.redirectRepository.Find(code)
	if err != nil {
		return redirect, err
	}

	redirect.incrementHits()
	err = rs.redirectRepository.Update(redirect)
	if err != nil {
		return redirect, err
	}

	return redirect, nil
}

// Store a redirect and return a unique code
func (rs *redirectService) Store(redirect Redirect) (code string, err error) {
	if govalidator.IsURL(redirect.URL) == false {
		return code, ErrRedirectInvalid
	}

	redirect.Code = randomCode(8)
	redirect.Created = time.Now()
	redirect.Updated = time.Now()

	err = rs.redirectRepository.Store(redirect)
	if err != nil {
		return code, err
	}

	return redirect.Code, nil
}

func (rs *redirectService) Report() (redirects []Redirect, err error) {
	return rs.redirectRepository.FindAll()
}

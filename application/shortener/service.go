package shortener

// RedirectService lets you Find a Redirect by url or store a redirect
type RedirectService interface {
	Find(code string) (redirect Redirect, err error)
	Store(redirect Redirect) (code string, err error)
	Report() (redirects []Redirect, err error)
}

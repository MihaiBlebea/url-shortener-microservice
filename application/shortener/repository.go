package shortener

// RedirectRepository interface for storing and retrieving a redirect object from the persistence layer
type RedirectRepository interface {
	Find(code string) (redirect Redirect, err error)
	FindAll() (redirects []Redirect, err error)
	Store(redirect Redirect) (err error)
	Update(redirect Redirect) (err error)
}

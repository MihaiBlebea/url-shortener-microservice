package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MihaiBlebea/url-shortener/shortener"
	"github.com/go-redis/redis/v8"
)

// Repository implementation for redis
type Repository struct {
	client *redis.Client
}

// NewClient returns a new redis client
func NewClient(host, port string) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return client, err
	}
	return client, nil
}

// NewRepository returns a new redis repository
func NewRepository(client *redis.Client) (repository *Repository) {
	return &Repository{client}
}

// Find method receives a code string and returns a redirect struct
func (r *Repository) Find(code string) (redirect shortener.Redirect, err error) {
	key := createKey(code)
	ctx := context.Background()

	// if r.client.HExists(ctx, key).Result()
	fields, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return redirect, err
	}
	redirect.URL = fields["url"]
	redirect.Code = fields["code"]
	hits, err := strconv.Atoi(fields["hits"])
	if err != nil {
		return redirect, err
	}
	redirect.Hits = hits

	created, err := strconv.Atoi(fields["created"])
	redirect.Created = time.Unix(int64(created), 0)

	updated, err := strconv.Atoi(fields["updated"])
	redirect.Updated = time.Unix(int64(updated), 0)

	return redirect, nil
}

// FindAll returns all redirects from redis
func (r *Repository) FindAll() (redirects []shortener.Redirect, err error) {
	return redirects, nil
}

// Store saves a redirect in redis
func (r *Repository) Store(redirect shortener.Redirect) (err error) {
	key := createKey(redirect.Code)
	ctx := context.Background()

	r.client.HSet(ctx, key, "url", redirect.URL)
	r.client.HSet(ctx, key, "code", redirect.Code)
	r.client.HSet(ctx, key, "hits", redirect.Hits)
	r.client.HSet(ctx, key, "created", redirect.Created.Unix())
	r.client.HSet(ctx, key, "updated", redirect.Updated.Unix())

	return nil
}

// Update receives a redirect struct and updates it's fields in redis
func (r *Repository) Update(redirect shortener.Redirect) (err error) {
	key := createKey(redirect.Code)
	ctx := context.Background()

	r.client.HIncrBy(ctx, key, "hits", 1)
	r.client.HSet(ctx, key, "updated", time.Now().Unix())

	return nil
}

func createKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

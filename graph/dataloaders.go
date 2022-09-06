package graph

import (
	"context"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/kazukiii/gqlgen-todos/graph/model"
	"net/http"
)

type Loaders struct {
	UserById *dataloader.Loader
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserById: dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
				fmt.Println("batch get users:", keys)
				userIds := keys.Keys()
				results := make([]*dataloader.Result, len(userIds))
				for i, id := range userIds {
					results[i] = &dataloader.Result{
						Data:  &model.User{ID: id, Name: "user " + id},
						Error: nil,
					}
				}
				return results
			}),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

type contextKey int

var loadersKey contextKey

func ctxLoaders(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

package graph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/msdsm/gqlgen-todos/graph/model"
)

type Loaders struct {
	// 実際は複数のdataloaderをここに集約
	UserById *dataloader.Loader[string, *model.User]
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserById: dataloader.NewBatchedLoader(func(ctx context.Context, userIds []string) []*dataloader.Result[*model.User] {
				fmt.Println("batch get users:", userIds)
				// ユーザーIDのリストから一括でユーザー情報を取得
				// 実際のアプリケーションではここでデータベースやAPIから一括取得する
				// 例: SELECT * FROM users WHERE id IN (userIds)
				results := make([]*dataloader.Result[*model.User], len(userIds))

				// ユーザー情報をマップに格納して高速にアクセスできるようにする
				userMap := make(map[string]*model.User)
				for _, id := range userIds {
					// 実際のアプリケーションではここでデータベースやAPIから取得したデータを使用
					userMap[id] = &model.User{
						ID:   id,
						Name: "user " + id,
					}
				}

				// 結果を組み立て
				for i, id := range userIds {
					if user, ok := userMap[id]; ok {
						results[i] = &dataloader.Result[*model.User]{
							Data:  user,
							Error: nil,
						}
					} else {
						results[i] = &dataloader.Result[*model.User]{
							Data:  nil,
							Error: fmt.Errorf("user not found: %s", id),
						}
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

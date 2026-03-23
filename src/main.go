package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-firestore/internal/infrastructure/database"
	"go-firestore/internal/infrastructure/persistence/firestore/message"
)

// メッセージ作成用リクエストボディの構造体
type CreateMessagesRequestBody struct {
	UID  string `json:"uid"`
	Text string `json:"text"`
}

// メッセージ更新用リクエストボディの構造体
type UpdateMessagesRequestBody struct {
	SenderID string `json:"senderId"`
	Text     string `json:"text"`
}

// メッセージのレスポンス結果用の構造体
type MessageResponse struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"senderId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func main() {
	// echoのルーター設定
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		// レスポンス結果の設定
		res := map[string]string{
			"message": "Hello World !!",
		}

		return c.JSON(http.StatusOK, res)
	})

	// サンプルAPI（CRUD処理）を追加
	apiV1 := e.Group("/api/v1")

	// メッセージ作成
	apiV1.POST("/messages", func(c echo.Context) error {
		// リクエストボディの取得
		var reqBody CreateMessagesRequestBody
		if err := c.Bind(&reqBody); err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Firestoreクライアント取得
		client, err := database.NewFirestoreClient(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create Firestore client: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}
		defer client.Close()

		// ドキュメント作成
		docRef, _, err := client.Collection("messages").Add(ctx, map[string]interface{}{
			"senderId":  reqBody.UID,
			"text":      reqBody.Text,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		})
		if err != nil {
			errMsg := fmt.Sprintf("failed adding a new message: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		// ドキュメント結果の取得
		docSnap, err := docRef.Get(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get a new message: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		// マッピング
		var msgDoc message.MessageDocument
		if err := docSnap.DataTo(&msgDoc); err != nil {
			errMsg := fmt.Sprintf("failed to DataTo struct: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		// レスポンス結果の設定
		message := MessageResponse{
			ID:        docSnap.Ref.ID,
			SenderID:  msgDoc.SenderID,
			Text:      msgDoc.Text,
			CreatedAt: msgDoc.CreatedAt,
			UpdatedAt: msgDoc.UpdatedAt,
		}

		return c.JSON(http.StatusCreated, message)
	})

	// 全てのメッセージ取得
	apiV1.GET("/messages", func(c echo.Context) error {
		ctx := c.Request().Context()

		// Firestoreクライアント取得
		client, err := database.NewFirestoreClient(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create Firestore client: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}
		defer client.Close()

		// イテレーター取得
		iter := client.Collection("messages").Documents(ctx)
		defer iter.Stop()

		// メッセージ取得処理
		var messages []MessageResponse

		for {
			// ドキュメント取得
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				errMsg := fmt.Sprintf("failed to iterate the list of messages: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
			}

			// マッピング
			var msgDoc message.MessageDocument
			if err := doc.DataTo(&msgDoc); err != nil {
				errMsg := fmt.Sprintf("failed to DataTo struct: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
			}

			// レスポンス結果の設定
			msgRes := MessageResponse{
				ID:        doc.Ref.ID,
				SenderID:  msgDoc.SenderID,
				Text:      msgDoc.Text,
				CreatedAt: msgDoc.CreatedAt,
				UpdatedAt: msgDoc.UpdatedAt,
			}

			messages = append(messages, msgRes)
		}

		// データが0件の場合、空の配列を設定
		if len(messages) == 0 {
			messages = []MessageResponse{}
		}

		return c.JSON(http.StatusOK, messages)
	})

	// 対象のメッセージ取得
	apiV1.GET("/messages/:id", func(c echo.Context) error {
		// リクエストパラメータの取得
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "id is required")
		}

		ctx := c.Request().Context()

		// Firestoreクライアント取得
		client, err := database.NewFirestoreClient(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create Firestore client: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}
		defer client.Close()

		// 対象データ取得
		docSnap, err := client.Collection("messages").Doc(id).Get(ctx)
		if err != nil {
			// 対象データが存在しない場合は空のオブジェクトを返す
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				return c.JSON(http.StatusOK, map[string]interface{}{})
			} else {
				errMsg := fmt.Sprintf("failed to get a message: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
			}
		}

		// マッピング
		var msgDoc message.MessageDocument
		if err := docSnap.DataTo(&msgDoc); err != nil {
			errMsg := fmt.Sprintf("failed to DataTo struct: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		// レスポンス結果の設定
		message := MessageResponse{
			ID:        docSnap.Ref.ID,
			SenderID:  msgDoc.SenderID,
			Text:      msgDoc.Text,
			CreatedAt: msgDoc.CreatedAt,
			UpdatedAt: msgDoc.UpdatedAt,
		}

		return c.JSON(http.StatusOK, message)
	})

	// 対象メッセージ更新
	apiV1.PUT("/messages/:id", func(c echo.Context) error {
		// リクエストパラメータの取得
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "id is required")
		}

		// リクエストボディの取得
		var reqBody UpdateMessagesRequestBody
		if err := c.Bind(&reqBody); err != nil {
			return err
		}

		ctx := c.Request().Context()

		// Firestoreクライアント取得
		client, err := database.NewFirestoreClient(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create Firestore client: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}
		defer client.Close()

		// 更新処理（トランザクション利用）
		docRef := client.Collection("messages").Doc(id)

		err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			docSnap, err := tx.Get(docRef)
			if err != nil {
				// 対象データが存在しない場合は404エラー
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.NotFound {
					return echo.NewHTTPError(http.StatusNotFound, "message not found")
				} else {
					errMsg := fmt.Sprintf("failed to get a message: %v", err)
					return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
				}
			}

			// マッピング
			var msgDoc message.MessageDocument
			if err := docSnap.DataTo(&msgDoc); err != nil {
				errMsg := fmt.Sprintf("failed to DataTo struct: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
			}

			// senderIdの認可チェック
			if reqBody.SenderID != msgDoc.SenderID {
				return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			}

			// 更新処理
			return tx.Update(docRef, []firestore.Update{
				{Path: "text", Value: reqBody.Text},
				{Path: "updatedAt", Value: time.Now()},
			})
		})
		if err != nil {
			return err
		}

		// 更新データ取得
		docSnap, err := docRef.Get(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get a message: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		// マッピング
		var msgDoc message.MessageDocument
		if err := docSnap.DataTo(&msgDoc); err != nil {
			errMsg := fmt.Sprintf("failed to DataTo struct: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		// レスポンス結果の設定
		message := MessageResponse{
			ID:        docSnap.Ref.ID,
			SenderID:  msgDoc.SenderID,
			Text:      msgDoc.Text,
			CreatedAt: msgDoc.CreatedAt,
			UpdatedAt: msgDoc.UpdatedAt,
		}

		return c.JSON(http.StatusOK, message)
	})

	// 対象メッセージ削除
	apiV1.DELETE("/messages/:id", func(c echo.Context) error {
		// リクエストパラメータの取得
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "id is required")
		}

		ctx := c.Request().Context()

		// Firestoreクライアント取得
		client, err := database.NewFirestoreClient(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create Firestore client: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}
		defer client.Close()

		// 対象データの存在チェック
		_, err = client.Collection("messages").Doc(id).Get(ctx)
		if err != nil {
			// 対象データが存在しない場合は空のオブジェクトを返す
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				return echo.NewHTTPError(http.StatusNotFound, "message not found")
			} else {
				errMsg := fmt.Sprintf("failed to get a message: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
			}
		}

		// 対象データ削除
		_, err = client.Collection("messages").Doc(id).Delete(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to delete a message: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}

		return c.NoContent(http.StatusNoContent)
	})

	// ログ出力
	slog.Info("start go-firestore")

	// サーバー起動
	e.Logger.Fatal(e.Start(":8081"))
}

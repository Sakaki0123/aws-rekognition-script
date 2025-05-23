package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

func NewConfig(ctx context.Context, region string) aws.Config {
	awsCfg, err := awscfg.LoadDefaultConfig(
		ctx,
		awscfg.WithDefaultRegion(region),
	)
	if err != nil {
		panic(fmt.Sprintf("failed loading aws config, %v", err))
	}
	return awsCfg
}

func main() {
	ctx := context.Background()
	// 引数-urlから画像ファイルのURLを取得
	fileURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cf/Flag_of_the_NSDAP_%281920%E2%80%931945%29.svg/250px-Flag_of_the_NSDAP_%281920%E2%80%931945%29.svg.png"

	// 画像ファイルを取得
	image, err := http.Get(fileURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer image.Body.Close()

	// 画像ファイルのデータを全て読み込み
	bytes, err := io.ReadAll(image.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// セッション作成
	cfg := NewConfig(ctx, "ap-northeast-1")
	svc := rekognition.NewFromConfig(cfg)
	// Rekognitionクライアントを作成

	// DetectFacesに渡すパラメータを設定
	params := &rekognition.DetectModerationLabelsInput{
		Image: &types.Image{
			Bytes: bytes,
		},
		// Confidenceが70%以上
		MinConfidence: aws.Float32(0.5),
	}

	// DetectFacesを実行
	resp, err := svc.DetectModerationLabels(ctx, params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	// 結果を出力
	// fmt.Printf("%+v\n", resp.ModerationLabels)
	// fmt.Printf("%+v\n", resp.ResultMetadata)
	for _, label := range resp.ModerationLabels {
		fmt.Printf("Name: %s, Confidence: %f\n", *label.Name, *label.Confidence)
	}
}

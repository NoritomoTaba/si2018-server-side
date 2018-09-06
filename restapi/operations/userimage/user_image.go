package userimage

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	
	"os"
	//"image"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	i := repositories.NewUserImageRepository()
	t := repositories.NewUserTokenRepository()
	
	// Params から, Image,Tokenの取得
	image := p.Params.Image
	token := p.Params.Token
	
	// 取得したtokenからUserTokenの取得
	loginUserToken,err := t.GetByToken(token)
	if err != nil {
		outPutPostStatus(500)
	}
	if loginUserToken == nil {
		outPutPostStatus(401)
	}
	
	// プロフィール画像の保存先情報の取得
	// local
	localpath := os.Getenv("ASSETS_PATH")
	// db
	dbpath := os.Getenv("ASSETS_BASE_URI")
	
	// プロフィール画像の保存する際の名前
	var fileName string
	fileName = token + ".png"
	
	// localに保存するプロフィール画像のPATH
	var filePathLocal string
	filePathLocal = localpath + fileName
	// dbに保存するプロフィール画像のPATH
	var filePathDB string
	filePathDB = dbpath + fileName
	
	/*
	// 取得したImageの形式を調べる
	_, format , err := image.DecodeConfig(image)
	if err != nil {
		fmt.Print("error")
	}
	fmt.Print(format)
	// プロフィール画像の名前
	// 拡張子によって名前を変更
	switch format {
	case "png":  fileName += ".png"
	case "jpeg": fileName += ".jpeg"
	case "gif":  fileName += ".gif"
	}
	*/
	
	// localに空のimage fileを用意
	file , _ := os.Create(filePathLocal)
	defer file.Close()
	file.Write(image)
	
	var updatedProfile entities.UserImage
	updatedProfile.UserID    = loginUserToken.UserID
	updatedProfile.Path      = filePathDB
	
	err = i.Update(updatedProfile)
	
	if err != nil {
		outPutPostStatus(500)
	}
	
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(filePathDB),
		})
}

func outPutPostStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	case 400:
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return nil
}
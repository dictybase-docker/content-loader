package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dictyBase/go-genproto/dictybaseapis/api/jsonapi"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/go-genproto/dictybaseapis/user"
	minio "github.com/minio/minio-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/urfave/cli.v1"
)

func importAction(c *cli.Context) error {
	logger := getLogger(c)
	s3Client, err := getS3Client(c)
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", c.GlobalString("content-api-host"), c.GlobalString("content-api-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to connect to content-api grpc backend %s", err),
			2,
		)
	}
	defer conn.Close()
	client := pb.NewContentServiceClient(conn)

	uconn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", c.String("user-api-host"), c.String("user-api-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to connect to user-api grpc backend %s", err),
			2,
		)
	}
	defer uconn.Close()
	uclient := user.NewUserServiceClient(uconn)
	user, err := uclient.GetUserByEmail(
		context.Background(),
		&jsonapi.GetEmailRequest{Email: c.String("email")},
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in retrieving user %s %s", c.String("email"), err),
			2,
		)
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := s3Client.ListObjects(c.String("s3-bucket"), c.String("remote-path"), true, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			return cli.NewExitError(
				fmt.Sprintf("unable to fetch object %s", object.Err),
				2,
			)
		}
		reader, err := s3Client.GetObject(c.String("s3-bucket"), object.Key, minio.GetObjectOptions{})
		if err != nil {
			logger.Errorf("unable to retrieve the object %s %s", object.Key, err)
			continue
		}
		defer reader.Close()
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			logger.Errorf("error in reading content for object %s %s", object.Key, err)
			continue
		}

		fsl := strings.Split(filepath.Base(object.Key), ".")
		if len(fsl) <= 1 && fsl[1] != c.String("extension") {
			logger.Warnf("skipping object %s", object.Key)
			continue
		}
		slug := fmt.Sprintf("%s-%s", c.String("namespace"), fsl[0])
		ct, err := client.GetContentBySlug(
			context.Background(),
			&pb.ContentRequest{Slug: slug},
		)
		if err != nil {
			st, ok := status.FromError(err)
			// insert new record
			if ok && st.Code() == codes.NotFound {
				logger.Warnf("no slug %s found, need to create record", slug)
				nc, err := client.StoreContent(context.Background(), &pb.StoreContentRequest{
					Data: &pb.StoreContentRequest_Data{
						Type: "contents",
						Attributes: &pb.NewContentAttributes{
							Name:      fsl[0],
							CreatedBy: user.Data.Id,
							Content:   string(b),
							Namespace: c.String("namespace"),
						},
					},
				})
				if err != nil {
					logger.Errorf("error in creating content for file %s %s", fsl[0], err)
					continue
				}
				logger.Infof(
					"created new content for file %s with id %d and slug %s",
					fsl[0],
					nc.Data.Id,
					nc.Data.Attributes.Slug,
				)
			} else {
				logger.Warnf("unknown error in finding slug %s %s", slug, st.Message())
			}
			continue
		}
		// update record
		uct, err := client.UpdateContent(
			context.Background(),
			&pb.UpdateContentRequest{
				Id: ct.Data.Id,
				Data: &pb.UpdateContentRequest_Data{
					Type: ct.Data.Type,
					Id:   ct.Data.Id,
					Attributes: &pb.ExistingContentAttributes{
						UpdatedBy: user.Data.Id,
						Content:   string(b),
					},
				},
			})
		if err != nil {
			logger.Errorf(
				"error in updating content for file %s slug %s %s",
				fsl[0], ct.Data.Attributes.Slug,
				err,
			)
			continue
		}
		logger.Infof(
			"update content for file %s with id %d and slug %s",
			fsl[0],
			uct.Data.Id,
			uct.Data.Attributes.Slug,
		)
	}
	return nil
}

func getS3Client(c *cli.Context) (*minio.Client, error) {
	s3Client, err := minio.New(
		fmt.Sprintf("%s:%s", c.String("s3-host"), c.String("s3-port")),
		c.String("access-key"),
		c.String("secret-key"),
		false,
	)
	if err != nil {
		return s3Client, fmt.Errorf("unable create the client %s", err.Error())
	}
	return s3Client, nil
}

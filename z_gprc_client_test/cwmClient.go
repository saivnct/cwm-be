package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"sol.go/cwm/fileHandler"
	"sol.go/cwm/proto/cwmSignalMsgPb"
	"sol.go/cwm/proto/grpcCWMPb"
	"sol.go/cwm/utils"
	"strconv"
)

func main() {
	fmt.Println("Hello, I'm a client")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	tls, _ := strconv.ParseBool(os.Getenv("GRPC_TLS"))
	fmt.Println("tls", tls)

	opts := []grpc.DialOption{}

	if tls {
		certFile := "../ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	//opts = append(opts, grpc.WithPerRPCCredentials(&appgrpc.GrpcLoginCreds{
	//	Authorization: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI4NDk2NjIyOTI2OCIsImV4cCI6MTY1NzcwNzUxM30.hFvlDngLDlArySCH1oyYfAAgbTB4cX8N2VsMM0EiQ6Q",
	//}))

	cc, err := grpc.Dial("localhost:"+os.Getenv("GRPC_PORT"), opts...)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer cc.Close() //this will call at very end of code
	c := grpcCWMPb.NewCWMServiceClient(cc)

	//creatUser(&c)
	//verifyAuthenCode(&c)
	//login(&c)
	//uploadFile(&c)
	downloadFile(&c)

	log.Println("Done!")
}

func creatUser(c *grpcCWMPb.CWMServiceClient) {
	request := grpcCWMPb.CreatAccountRequest{
		Phone:       "966229268",
		CountryCode: "84",
	}

	response, err := (*c).CreatUser(context.Background(), &request)
	if err != nil {

		log.Fatalf("could not send CreatUser Request: %v - %v", status.Code(err), err)
	}

	spew.Dump(response)
}

func verifyAuthenCode(c *grpcCWMPb.CWMServiceClient) {
	deviceInfo := grpcCWMPb.DeviceInfo{
		DeviceName:   "S22",
		Imei:         "123456789",
		Manufacturer: "Samsung",
		Os:           "Android",
		OsVersion:    "26",
	}

	request := grpcCWMPb.VerifyAuthencodeRequest{
		Phone:       "966229268",
		CountryCode: "84",
		DeviceInfo:  &deviceInfo,
		Authencode:  "058835",
	}

	response, err := (*c).VerifyAuthencode(context.Background(), &request)
	if err != nil {
		log.Fatalf("could not send VerifyAuthencode Request: %v", err)
	}

	spew.Dump(response)

}

func login(c *grpcCWMPb.CWMServiceClient) {
	password := "0f8a59c9-21ed-4539-b37b-5757db02cc98"
	nonce := "fe96b570-15ca-4358-a2cf-674023e98124"
	nonceResponse := utils.GetNonceRespone(password, nonce)

	request := grpcCWMPb.LoginRequest{
		PhoneFull: "84966229268",
		Nonce:     nonce,
		Response:  nonceResponse,
	}

	response, err := (*c).Login(context.Background(), &request)
	if err != nil {
		//NOT SUPPORT FOR JAVA PROTOBUF_LITE
		//st := status.Convert(err)
		//if st.Code() == codes.Unauthenticated {
		//	for _, detail := range st.Details() {
		//		switch t := detail.(type) {
		//		case *errdetails.BadRequest:
		//			for _, violation := range t.GetFieldViolations() {
		//				spew.Dump(violation)
		//			}
		//		}
		//	}
		//
		//	log.Println("Unauthenticated - response", st.Details())
		//}
		//log.Fatalf("could not send Login Request: %v - %v", status.Code(err), err)
		log.Fatalf("could not send Login Request: %v", status.Convert(err).Message())
	}

	spew.Dump(response)
}

func uploadFile(c *grpcCWMPb.CWMServiceClient) {
	imagePath := "/Users/solgo/Downloads/chat-platform.xlsx"
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open file: ", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("cannot obtain file stat: ", err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	stream, err := (*c).UploadMediaMsg(context.Background())
	if err != nil {
		log.Fatal("cannot upload file: ", err)
	}

	log.Printf("file size: %d, checksum %v", fileInfo.Size(), checksum)

	req := &grpcCWMPb.UploadMediaMsgRequest{
		Data: &grpcCWMPb.UploadMediaMsgRequest_MediaMsgInfo{
			MediaMsgInfo: &grpcCWMPb.MediaMsgInfo{
				MsgId:     "123456789",
				MediaType: cwmSignalMsgPb.SIGNAL_MEDIA_TYPE_DOC,
				Checksum:  checksum,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send file info to server: ", err, stream.RecvMsg(nil))
	}

	file.Seek(0, io.SeekStart)
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &grpcCWMPb.UploadMediaMsgRequest{
			Data: &grpcCWMPb.UploadMediaMsgRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			//To get the real error that contains the gRPC status code, we must call stream.RecvMsg() with a nil parameter. The nil parameter basically means that we don't expect to receive any message, but we just want to get the error that function returns
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("file uploaded with id: %s, name: %s, msgId: %s", res.GetFileId(), res.GetFileName(), res.GetMsgId())
}

func downloadFile(c *grpcCWMPb.CWMServiceClient) {

	req := &grpcCWMPb.DownloadMediaMsgRequest{
		FileId: "2b0888e0-d40f-4a97-b861-d83ee5a66c62",
		MsgId:  "d0b55f1f-4103-4e85-9b41-a6af100fcf63",
	}

	stream, err := (*c).DownloadMediaMsg(context.Background(), req)
	if err != nil {
		log.Fatal("cannot download file: ", err)
	}

	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		//log.Println("waiting to receive more media data")

		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("received file data from client")
			break
		}
		if err != nil {
			log.Fatal("Cannot receive chunk data: ", err)
		}

		chunk := res.GetChunkData()
		size := len(chunk)

		//log.Printf("received a chunk with size: %d\n", size)

		fileSize += size
		_, err = fileData.Write(chunk)
		if err != nil {
			log.Fatal("Cannot write chunk data: ", err)
		}
	}

	checksum := fmt.Sprintf("%x", md5.Sum(fileData.Bytes()))

	diskFileStore := fileHandler.NewDiskFileStore("media")
	_, err = diskFileStore.SaveToDisk("test", "jpg", fileData)

	log.Println("Done download file", fileSize, checksum)
	if err != nil {
		log.Fatal("cannot save image to the store: ", err)
	}
}
